<style>
td {
    white-space:nowrap;
    overflow:hidden;
    text-overflow: ellipsis;
    font-family: Consolas,Monaco,monospace;
}
ul.pagination {
    display: inline-block;
    padding: 0;
    margin: 0 auto;
text-align:center;
}

ul.pagination li {display: inline;}

ul.pagination li a {
    color: black;
    float: left;
    padding: 8px 16px;
    text-decoration: none;
    border-radius: 15px;
}
ul.pagination li a.active {
    background-color: #4CAF50;
    color: white;
    border-radius: 15px;
}

ul.pagination li a:hover:not(.active) {background-color: #ddd;}
</style>
<center>
<ul class="pagination" >
  <li><a href="/view/blocks/{{.Page.FirstPage}}">Fir</a></li>
  <li><a href="/view/blocks/{{.Page.PrevPage}}">Pre </a></li>
  <li><a class="active" href="#">{{.Page.CurrentPage}} of {{.Page.LastPage}}</a></li>
  <li><a href="/view/blocks/{{.Page.NextPage}}">Next</a></li>
  <li><a href="/view/blocks/{{.Page.LastPage}}">Last</a></li>
</ul>
<!--  共{{.Page.Items}}条记录 -->
 </center>

<div class="row">
        <!-- Begin: life time stats -->
        <div class="portlet light portlet-fit portlet-datatable ">
            <div class="table-scrollable table-scrollable-borderless">
                    <table class="table table-hover table-light" >
                        <thead>
                            <tr class="uppercase">
                                <th> Block Hash</th>
                                  <th> Height </th>
                                <th> Validator </th>
                                <th> Txs </th>
                                <th> Time </th>
                                <th> Interval(-25) </th>
                            </tr>
                        </thead>
                        <tbody id="dashboard-block-table">
                            {{range $index,$block := .Blocks}}
                                  <tr>
                                    <td title={{$block.Hash}}>
                                        <a href="/view/blocks/hash/{{$block.Hash}}">{{$block.Hash}}</a>
                                    </td>
                                    <td>{{$block.Height}}</td>
                                    <td>{{$block.ValidatorsHash}}</td>
                                    <td>{{$block.NumTxs}}</td>
                                    <td>{{$block.Time}}</td>
                                    <td>{{$block.Interval}}</td>
                                  </tr>
                             {{end}}
                        
                        </tbody>
                    </table>
            </div>
</div>
