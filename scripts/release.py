#!/usr/bin/env python3
import hashlib, json, os, re, shutil, struct, subprocess, zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
VERSION = '0.1.0'
GO_VERSION = 'go1.26.8'
TOKEN_PATTERN = re.compile(rb'(?<![A-Za-z0-9])\d{6,}:[A-Za-z0-9_-]{20,}(?![A-Za-z0-9])')
PRIVATE_ID_PATTERN = re.compile(rb'"(?:telegram_)?(?:chat_id|channel_id|username)"\s*:\s*"(?:-?[0-9]{4,}|@[A-Za-z0-9_]+)"')
STAMP = (2026, 10, 6, 0, 0, 0)
go = os.environ.get('GO', 'go')
env = dict(os.environ, CGO_ENABLED='0', GOOS='linux', GOARCH='arm', GOARM='5', GOTOOLCHAIN='local')
actual = subprocess.check_output([go, 'version'], text=True).split()[2]
if actual != GO_VERSION:
    raise SystemExit(f'Requires {GO_VERSION}; got {actual}')

dist = ROOT / 'dist'
dist.mkdir(exist_ok=True)
stage = ROOT / 'build/release'
if stage.exists(): shutil.rmtree(stage)
app = stage / '.adds/nickelgram'
(app / 'lib').mkdir(parents=True)
(stage / '.adds/nm').mkdir(parents=True)

allowed = ['nickelgram.sh','config.example.json','cacert.pem','LICENSE-Go','LICENSE-MPL-2.0','sqlite3','lib/libsqlite3.so.0']
for name in allowed:
    dst = app / name
    dst.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(ROOT / 'payload/.adds/nickelgram' / name, dst)
shutil.copy2(ROOT / 'payload/.adds/nm/nickelgram', stage / '.adds/nm/nickelgram')

subprocess.run([go,'build','-trimpath','-buildvcs=false','-ldflags=-s -w -buildid=', '-o',str(app/'nickelgram'),'./cmd/nickelgram'], cwd=ROOT, env=env, check=True)

binary=(app/'nickelgram').read_bytes()
assert binary[:6] == b'\x7fELF\x01\x01'
assert struct.unpack_from('<H',binary,18)[0] == 40
phoff=struct.unpack_from('<I',binary,28)[0]
entsize,count=struct.unpack_from('<HH',binary,42)
assert all(struct.unpack_from('<I',binary,phoff+i*entsize)[0] not in (2,3) for i in range(count))

for source,target in [('README.md','README.zh-CN.md'),('LICENSE','LICENSE'),('docs/THIRD_PARTY.md','docs/THIRD_PARTY.md')]:
    (app/target).parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(ROOT/source, app/target)

manifest={'version':VERSION,'go':actual,'target':'linux/arm GOARM=5 CGO_ENABLED=0','hardware_tested':False,
          'files':{p.relative_to(stage).as_posix():hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted(stage.rglob('*')) if p.is_file()}}
(app/'MANIFEST.json').write_text(json.dumps(manifest,indent=2,ensure_ascii=False)+'\n')

def archive(path, files):
    with zipfile.ZipFile(path,'w',zipfile.ZIP_DEFLATED,compresslevel=9) as z:
        for src,name,executable in sorted(files,key=lambda x:x[1]):
            info=zipfile.ZipInfo(name,STAMP); info.create_system=3
            info.external_attr=(0o100755 if executable else 0o100644)<<16
            info.compress_type=zipfile.ZIP_DEFLATED
            z.writestr(info,src.read_bytes())
    with zipfile.ZipFile(path) as z:
        assert z.testzip() is None
        for n in z.namelist():
            assert not n.endswith('/config.json') and '/state/' not in n
            data=z.read(n)
            assert not TOKEN_PATTERN.search(data), f'Bot token pattern found in {n}'
            assert not PRIVATE_ID_PATTERN.search(data), f'Private Telegram ID found in {n}'

package=dist/f'NickelGram-{VERSION}-kobo-arm.zip'
archive(package,[(p,p.relative_to(stage).as_posix(),p.name in ('nickelgram.sh','nickelgram','sqlite3')) for p in stage.rglob('*') if p.is_file()])

source_files=[ROOT/n for n in ('.gitignore','go.mod','LICENSE','README.md','test.sh')]
for folder in ('cmd','scripts','docs','payload'):
    source_files += [p for p in (ROOT/folder).rglob('*') if p.is_file() and '__pycache__' not in p.parts and p.suffix != '.pyc']
assert not any(p.name=='config.json' or 'state' in p.relative_to(ROOT).parts for p in source_files)
for p in source_files:
    data=p.read_bytes()
    if p.name=='config.json': raise SystemExit('Refusing to package config.json')
    if TOKEN_PATTERN.search(data) or PRIVATE_ID_PATTERN.search(data): raise SystemExit(f'Possible Telegram credential in {p}')

source_zip=dist/f'NickelGram-{VERSION}-source.zip'
archive(source_zip,[(p,f'NickelGram-{VERSION}/'+p.relative_to(ROOT).as_posix(),p.suffix=='.sh') for p in source_files])
(dist/'SHA256SUMS').write_text(''.join(hashlib.sha256(p.read_bytes()).hexdigest()+'  '+p.name+'\n' for p in (package,source_zip)))
print(package)
print(source_zip)
