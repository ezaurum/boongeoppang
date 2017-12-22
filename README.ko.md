[English](README.en.md)
# 붕어빵

관습에 따라 템플릿 파일을 불러오고, `string`으로 된 키를 이용해서 템플릿에 접근할 수있도록 해 준다.

## 사용법

``` 
	container := boongeoppang.Load("tests/full")
```

## 디렉토리 구조
템를릿 디렉토리 안에 `_default`,`_partials`가 있으면 기본 레이아웃과 

`user/profile` 에 대해 템플릿 파일 찾는 순서
``` 
/_default/baseof
/_default/profile
/user/profile
```
없는 파일은 불러오지 않는다.

```
_default/ - 전체 레이아웃이 들어가 있는 디렉토리
    baseof.tmpl - 기본 전체 레이아웃 틀. 없으면 로드 안함
    
    single.tmpl - single layout 형태 기본 
    test.tmpl
    
partials/
    모든 파일 기본 로드

content/
    single.tmpl -  "content", single 을 로드할 때 로드
user/
    test.tmpl
    
```

