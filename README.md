# White Walker

컨벤션에 따라 템플릿 파일 로드할 수 있도록 경로 준비

```
_default/
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

