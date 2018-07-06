<style>
    td {
          white-space:nowrap;
          overflow:hidden;
          text-overflow: ellipsis;
  font-family: Consolas,Monaco,monospace;
}
</style>

<div class="row">
        <!-- Begin: life time stats -->
        <div class="portlet light portlet-fit portlet-datatable ">
            <div class="portlet-title">
                <div class="caption">
                    <i class="fa fa-book font-green"></i>
                    <span class="caption-subject font-green sbold" id="last_block"></span>
                </div>
                <div class="actions">
                </div>
            </div>
                 <div class="portlet-title">
                    <div class="caption">
                        <i class="fa fa-book font-green"></i>
                        <span class="caption-subject font-green sbold" id="txs"></span>
                    </div>
                    <div class="actions">
                    </div>
                </div>

        </div>
        <!-- End: life time stats -->
    </div>
</div>

<style>
    td {
          white-space:nowrap;
          overflow:hidden;
          text-overflow: ellipsis;
    }
</style>


<script>
    $.extend( $.fn.dataTable.defaults, {
        searching: false,
        ordering:  false
    });

    var BlockDashboardInit = function() {
        var renderBlockList = function(data) {
            $('#last_block').text("Latest Block Height: "+ data.LatestBlockHeight);
            $('#txs').text("Transactions: " + data.Transactions);
        }
        var flushBlockList = function() {
            $.ajax({
                url: '/view/',
                type: 'GET',
                contentType: "application/json; charset=utf-8",
                dataType: "json",
                success: function(result) {
                    console.log(result)
                    if (result.success) {
                        renderBlockList(result.data)
                        setTimeout(function() {
                            flushBlockList();
                        }, 5000);
                    }
                }
            });
        }

        flushBlockList();
    };

    $(function() {
        BlockDashboardInit();
    });
</script>

