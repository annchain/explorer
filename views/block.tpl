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
                    <span class="caption-subject font-green sbold">Hash : {{.Block.Hash}}</span>
                </div>
                <div class="actions">
                </div>
            </div>
            <div class="portlet-body">
                <div class="table-container">
                    <table class="table table-striped table-bordered table-hover">
                        <tbody>
                          <tr><td>Height:</td><td>{{.Block.Height}}</td></tr>
                          <tr><td>NumTxs:</td><td>{{.Block.NumTxs}}</td></tr>
                          <tr><td>ValidatorsHash:</td><td>{{.Block.ValidatorsHash}}</td></tr>
                          <tr><td>AppHash:</td><td>{{.Block.AppHash}}</td></tr>
                          <tr><td>Reward:</td><td>{{.Block.Reward}}&nbsp;CHO</td></tr>
                          <tr><td>Winner:</td><td>{{.Block.CoinBase}}</td></tr>
                          <tr><td>Time:</td><td>{{.Block.Time}}</td></tr>
                        </tbody>
                    </table>
                </div>

                 <div class="portlet-title">
                    <div class="caption">
                        <i class="fa fa-book font-green"></i>
                        <span class="caption-subject font-green sbold">Transactions : </span>
                    </div>
                    <div class="actions">
                    </div>
                </div>
                <div class="portlet-body">
                    <div class="table-container">
                        <table class="table table-striped table-bordered table-hover">
                            <tbody>
                            {{range $index,$tx := .Transactions}}
                                  <tr><td>Hash:</td><td>{{$tx.Hash}}</td></tr>
                                  <tr>
					<td>Payload:</td>
					<td>
						<textarea style="width:100%;heigth:100px"> {{$tx.PayloadHex}}</textarea>
					</td>
				 </tr>
                             {{end}}
                            </tbody>
                        </table>
                    </div>
                </div>

            </div>
        </div>
        <!-- End: life time stats -->
    </div>
</div>
