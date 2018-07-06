<style>
    td {
          white-space:nowrap;
          overflow:hidden;
          text-overflow: ellipsis;
    }
</style>

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
                        var url = window.location.href
                        var hash = url.substring(url.lastIndexOf("\/")+1,url.lenght)
                        window.location.href = "/view/h5/txs/hash/"+hash

                    } else {
                        // pc端

                        console.log("pc")
                    }
    </script>

<div class="row">
    <div class="col-md-12">
        <!-- Begin: life time stats -->
        <div class="portlet light portlet-fit portlet-datatable ">
            <div class="portlet-title">
                <div class="caption">
                    <i class="fa fa-book font-green"></i>
                    <span class="caption-subject font-green sbold">Transaction : {{.Transaction.Hash}}</span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-container">
                    <table class="table table-striped table-bordered table-hover">
                        <tbody>

                          <tr><td>Block :</td><td><a href="/view/blocks/hash/{{.Transaction.Block}}">{{.Transaction.Block}}</td></tr>
                          <tr><td>Payload:</td><td><textarea style="width:100%;height:200px;" readonly> {{.Transaction.PayloadHex}} </textarea></td></tr>
                        </tbody>
                    </table>
                </div>



            </div>
        </div>
        <!-- End: life time stats -->
    </div>
</div>
