<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Annchain</title>
    <script>
      var browser = function () {
                        var u = navigator.userAgent, app = navigator.appVersion;
                        return {
                            //是否为移动终端
                            mobile: !!u.match(/AppleWebKit.*Mobile.*/),
                            //ios终端
                            ios: !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/),
                            //android终端
                            android: u.indexOf('Android') > -1 || u.indexOf('Adr') > -1,
                            //是否iPad
                            iPad: u.indexOf('iPad') > -1
                        };
                    }();

                    if(browser.mobile) {
                        // h5页面
                        console.log("mobile")
                    } else {
                        // pc端
                        var url = window.location.href
                        var hash = url.substring(url.lastIndexOf("\/")+1,url.lenght)

                        console.log("view/txs/hash"+hash)
                        window.location.href = "/view/txs/hash/"+hash
                        console.log("pc")
                    }
    </script>
    <style>
        html, body, p, div, h1, h2, ul, li {
            margin: 0;
            padding: 0;
            font-weight: normal;
        }
        ul, li {
            list-style: none;
        }
        body {
            background-color: #F8F8F8;
        }
        header {
            height: 1.2rem;
            font-size: .14rem;
            line-height: .2rem;
            padding-top: 1em;
            text-align: center;
            height: 1.2rem;
            color: #fff;
            box-sizing: border-box;
            background-image: linear-gradient(180deg, #4D6BC8 19%, #1D3ACE 92%);
        }
        .title {
            font-size: .18rem;
            line-height: .25rem;
            margin-bottom: .08rem;
            font-family: "PingFangSC-Semibold";
        }
        
        /* 页面中部height和numtx样式 */
        .height-numtx {
            position: absolute;
            left: .23rem;
            right: .23rem;
            top: .9rem;
            background-color: #fff;
            padding: 0 0.18rem;
            border-radius: .06rem;
            box-shadow: 0 .02rem .04rem 0 rgba(0, 0, 0, .13);
        }
        .height-numtx > li {
            display: flex;
            justify-content: space-between;
            align-items: center;
            line-height: 0.7rem;
            font-size: .14rem;
            color: #666;
        }
        .height-numtx > li > img {
            width: 0.18rem;
            height: 0.18rem;
            margin-right: 0.08rem;
        }
        .height-numtx > li > h2 {
            font-weight: bold;
            font-size: .14rem;
        }
        .height-numtx > li > h1 {
            flex: 1 1 0;
            font-size: .2rem;
            font-weight: bold;
            text-align: right;
        }
        .height-numtx > li:first-child {
            border-bottom: 1px solid #F0F0F0;
        }
        /* 页面中部列表样式 */
        .content {
            font-size: .12rem;
            line-height: .17rem;
            padding-top: 1.28rem;
            margin: 0 .4rem;
            color: #999;
        }
        .content .item {
            padding-bottom: .16rem;
            margin-bottom: .12rem;
            border-bottom: 1px solid #F0F0F0;
        }
        .content .item:last-child {
            margin-bottom: .04rem;
            border: none;
        }
        .content .item > h2 {
            font-weight: bold;
            color: #666;
            margin-bottom: .08rem;
            font-size: .12rem;
        }
        /* 页面底部Transactions样式 */
        .panel {
            padding: .25rem;
            color: #666;
            padding-bottom: .1rem;
            border-radius: .06rem;
            background-color: #fff;
            margin: 0 .22rem;
            box-shadow: 0 .02rem .04rem 0 rgba(0, 0, 0, .13);
        }
        .panel .title {
            margin-bottom: 0.24rem;
            text-align: center;
        }
        .panel h3 {
            font-size: .14rem;
            line-height: .2rem;
            margin-bottom: .06rem;
        }
        .panel > p {
            font-size: .12rem;
            line-height: 0.17rem;
            margin-bottom: .3rem;
            word-break: break-word;
        }
        /* 页面底部logo样式 */
        footer {
            font-size: 0;
            padding-top: 0.2rem;
            padding-bottom: 0.3rem;
            text-align: center;
        }
        footer > img {
            width: 1.12rem;
        }
    </style>
</head>
<body>
    <header>
        <h2 class="title">BlockHash</h2>
        <p>{{.Transaction.Block}}</p>
    </header>

    <ul class="height-numtx">
        <li>
            <img src="/assets/global/img/icon-height.png" alt="">
            <h2>Height</h2>
            <h1>{{.Transaction.Height}}</h1>
        </li>
        <li>
            <img src="/assets/global/img/icon-numtx.png" alt="">
            <h2>NumTxs</h2>
            <h1>{{.Block.NumTxs}}</h1>
        </li>
    </ul>

    <ul class="content">
        <li class="item">
            <h2>ValidatorsHash</h2>
            <p>{{.Block.ValidatorsHash}}</p>
        </li>
        <li class="item">
            <h2>AppHash</h2>
            <p>{{.Block.AppHash}}</p>
        </li>
        <li class="item">
            <h2>Time</h2>
            <p>{{.Transaction.Time}}</p>
        </li>
    </ul>
    <div class="panel">
        <h2 class="title">
            Transactions
        </h2>
        <h3>TxHash</h3>
        <p>{{.Transaction.Hash}}</p>
        <h3>Payload</h3>
        <p>
            {{.Transaction.PayloadHex}}
        </p>
    </div>
    <footer>
        <img src="/assets/global/img/logo.png" alt="">
    </footer>
    <script>



        function calcRem() {
            var winWidth = document.documentElement.clientWidth;
            var fontSize = 100 * winWidth / 375 + 'px';

            document.documentElement.style.fontSize = fontSize;
        }

        window.onresize = calcRem;

        calcRem();
    </script>
</body>
</html>