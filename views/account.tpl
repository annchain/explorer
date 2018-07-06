<style>
    td {
          white-space:nowrap;
          overflow:hidden;
          text-overflow: ellipsis;
  font-family: Consolas,Monaco,monospace;
}
</style>

<div class="row">
    <div class="col-md-12">
        <!-- Begin: life time stats -->
        <div class="portlet light portlet-fit portlet-datatable ">
            <div class="portlet-body">
                <div class="table-container">
                    <table class="table table-striped table-bordered table-hover">
                        <tbody>
                          <tr><td>Account Pubkey:</td><td>{{.Account.Pubkey}}</td></tr>
                          <tr><td>Balance:</td><td>{{.Account.Balance}}</td></tr>
                        </tbody>
                    </table>
                </div>


            </div>
        </div>
        <!-- End: life time stats -->
    </div>
</div>
