# gotpl
* go语言模板引擎。
* 通过”代码生成“方式，把模板文件(tpl/*.html), 生成为go源码(gen/*.go)，最后项目编译运行的都是go文件。
* 只需要掌握三个关键字：@extends,@block,@section,其他都是go语言的关键字。真的不能再简单了。

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

# 模板继承（@extends，@block）
base.tpl :
```
@{
    import()
    var val int
}

<html>
<body>
@block aa {
    <div>base aa content</div>

    @block bb {
        <div>"base bb content"</div>
    }

    @block cc {
        <div>base cc content</div>
    }
}

@section TestSection(val int)

</body>

</html>

```

child.tpl :

```
@extends base

@block bb {
    <div>child bb block content.</div>
}

@block cc {
    @for i:=0;i<10;i++ {
        <p>@i</p>
    }
}
```

语句 "@extends base", 意思是child.html继承base.html。模板继承方式，类似django，通过覆盖block。如覆盖base里面的bb，cc块。

# 模块组件（@section）
base.tpl :
```
@{
    import()
    var val int
}

<html>
<body>
@block aa {
    <div>base aa content</div>

    @block bb {
        <div>"base bb content"</div>
    }

    @block cc {
        <div>base cc content</div>
    }
}

@section TestSection(val int)

</body>

</html>

```

sections/test_section.tpl:
```
@section TestSection(val int) {
    <div>this is TestSection content. Param "val" is: @val </div>
}
```

section必须放在sections目录下，文件名不限制。

# LICENSE

[LICENSE](LICENSE)? Well, [WTFPL](http://www.wtfpl.net/about/).

# Todo
