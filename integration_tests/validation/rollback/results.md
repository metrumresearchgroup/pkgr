tags: rollback

result: PASS

date_run: 12-03-2019

## pkgr install will fail to install xml2, and after running, test-library will be identical to preinstalled-library
### before
![before](pil_before1.png)

### after is the same
![after](pil_after1.png)

## pkgr install --update will fail to install xml2, and after running, test-library will be identical to preinstalled-library
### before
![before](pil_before2.png)

### after is the same
![after](pil_after2.png)
