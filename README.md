# gotpl
go语言模板引擎。基于[gorazor](https://github.com/sipin/gorazor)开发。

# 特性
* 简洁优雅
* 模板继承
* 原生go语言
* 模块，组件

# 原生go语句

```
@if .... {
	....
}

@if .... {
	....
} else {
	....
}

@for .... {

}

@{
	switch .... {
	case ....:
	      <p>...</p>
	case 2:
	      <p>...</p>
	default:
	      <p>...</p>
	}
}
```

# 模板继承（extends，block）
base.tpl :
```
<html>

@block aa {
    aaaa

    @block bb {
        bbb
    }

    @block cc {
        ccc
    }
}

@section Pagination(curPage int)

</html>
```

test_extends_base.tpl :

```
@block bb {
     extends bbb
}

@block cc {
     @for i:=0;i<10;i++ {
     <p>@i</p>
     }
}
```

文件名test_extends_base.tpl,代表test.tpl继承base.tpl。模板继承方式，类似django，通过覆盖block。

# 模块组件（section）
base.tpl :
```
<html>

...

@section Pagination(curPage int)

</html>
```

sections/page.tpl:
```
@{
    import (

    )
}

@section Pagination(curPage int) {
    <div>curPage is: @curPage </div>
}
```

section必须放在sections目录下，文件名不限制。

# LICENSE

[LICENSE](LICENSE)? Well, [WTFPL](http://www.wtfpl.net/about/).

# Todo
