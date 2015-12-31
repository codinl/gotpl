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
@{
    import()
    var curPage int
}

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

child.tpl :

```
@extends base

@block bb {
    "this is content: extends bbb"
}

@block cc {
    @for i:=0;i<10;i++ {
        <p>@i</p>
    }
}
```

关键字 "extends base", child.tpl继承base.tpl。模板继承方式，类似django，通过覆盖block。如覆盖base里面的bb，cc块。

# 模块组件（section）
base.tpl :
```
@{
    import()
    var curPage int
}

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

sections/page.tpl:
```
@section Pagination(curPage int) {
    <div>curPage is: @curPage </div>
}
```

section必须放在sections目录下，文件名不限制。

# LICENSE

[LICENSE](LICENSE)? Well, [WTFPL](http://www.wtfpl.net/about/).

# Todo
