tags: multi-repo, repo-customizations, pkg-customizations, heavy

result: PASS

## Pkgr Plan equivalent to output in guide
![plan1](plan1.png)
![plan2](plan2.png)

Note: `pkgr plan` was run once before this output, which took care of some cacheing.
It looks like the guide did this as well, and that's not what's being tested here,
so this test is still valid.

## Pkgr Install
![install1](install1.png)
![install2](install2.png)
![install3](install3.png)
