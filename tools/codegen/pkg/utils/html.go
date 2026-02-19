package utils

type HTMLTemplateData struct {
	FeatureGates []HTMLFeatureGate
}

type HTMLFeatureGate struct {
	Name       string
	Sufficient bool
	Variants   []HTMLVariantColumn
	Tests      []HTMLTestRow
}

type HTMLVariantColumn struct {
	Topology     string
	Cloud        string
	Architecture string
	NetworkStack string
	ColIndex     int
}

type HTMLTestRow struct {
	TestName string
	Cells    []HTMLTestCell
}

type HTMLTestCell struct {
	PassPercent    int
	SuccessfulRuns int
	TotalRuns      int
	FailedRuns     int
	Failed         bool
}

const HTMLTemplateSrc = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>FeatureGate Promotion Summary</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap/4.6.1/css/bootstrap.min.css"
          integrity="sha512-T584yQ/tdRR5QwOpfvDfVQUidzfgc2339Lc8uBDtcp/wYu80d7jwBgAxbyMh0a9YM9F8N3tdErpFI8iaGx6x5g=="
          crossorigin="anonymous">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <style>
        body { padding: 20px; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
        .sortable { cursor: pointer; user-select: none; }
        .sortable:hover { background-color: #e9ecef; }
        .sort-indicator::after { content: ' \2195'; opacity: 0.3; font-size: 0.8em; }
        .sort-asc::after { content: ' \2191'; opacity: 1; font-size: 0.8em; }
        .sort-desc::after { content: ' \2193'; opacity: 1; font-size: 0.8em; }
        .fail-cell { background-color: #f8d7da; }
        .pass-cell { background-color: #d4edda; }
        .test-name { max-width: 500px; word-wrap: break-word; font-size: 0.85em; text-align: left; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 2em; }
        th, td { border: 1px solid #dee2e6; padding: 8px; text-align: center; vertical-align: middle; }
        th:first-child, td:first-child { text-align: left; }
        th { background-color: #f8f9fa; }
        .network-stack { font-weight: bold; color: #0056b3; }
        .alert { padding: 12px 20px; margin-bottom: 20px; border-radius: 4px; }
        .alert-success { background-color: #d4edda; border: 1px solid #c3e6cb; color: #155724; }
        .alert-danger { background-color: #f8d7da; border: 1px solid #f5c6cb; color: #721c24; }
        h1 { margin-bottom: 24px; }
        h2 { margin-top: 32px; margin-bottom: 12px; }
    </style>
</head>
<body>
<div class="container-fluid">
    <h1>FeatureGate Promotion Summary</h1>
    {{if not .FeatureGates}}<p>No new Default FeatureGates found.</p>{{end}}
    {{range $fgIdx, $fg := .FeatureGates}}
    <h2>{{$fg.Name}}</h2>
    {{if $fg.Sufficient}}
    <div class="alert alert-success">Sufficient CI testing for &quot;{{$fg.Name}}&quot;.</div>
    {{else}}
    <div class="alert alert-danger">
        <strong>INSUFFICIENT</strong> CI testing for &quot;{{$fg.Name}}&quot;.
        <ul>
            <li>At least five tests are expected for a feature</li>
            <li>Tests must be run on every TechPreview platform (ask for an exception if your feature doesn&#39;t support a variant)</li>
            <li>All tests must run at least 14 times on every platform</li>
            <li>All tests must pass at least 95% of the time</li>
        </ul>
    </div>
    {{end}}
    {{if $fg.Tests}}
    <table id="table-{{$fgIdx}}">
        <thead>
            <tr>
                <th class="sortable sort-indicator" data-table="{{$fgIdx}}" data-col="0" data-sort="text">Test Name</th>
                {{range $v := $fg.Variants}}
                <th class="sortable sort-indicator" data-table="{{$fgIdx}}" data-col="{{$v.ColIndex}}" data-sort="percent">
                    {{$v.Topology}}<br>{{$v.Cloud}}<br>{{$v.Architecture}}{{if $v.NetworkStack}}<br><span class="network-stack">{{$v.NetworkStack}}</span>{{end}}
                </th>
                {{end}}
            </tr>
        </thead>
        <tbody>
            {{range $test := $fg.Tests}}
            <tr>
                <td class="test-name">{{$test.TestName}}</td>
                {{range $cell := $test.Cells}}
                <td class="{{if $cell.Failed}}fail-cell{{else}}pass-cell{{end}}" data-pass-percent="{{$cell.PassPercent}}">
                    {{if $cell.Failed}}<strong>FAIL</strong><br>{{end}}
                    {{$cell.PassPercent}}% ({{$cell.SuccessfulRuns}} / {{$cell.TotalRuns}})
                    {{if gt $cell.FailedRuns 0}}<br><small>{{$cell.FailedRuns}} failed</small>{{end}}
                </td>
                {{end}}
            </tr>
            {{end}}
        </tbody>
    </table>
    {{end}}
    {{end}}
</div>
<script>
document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('th.sortable').forEach(function(th) {
        th.addEventListener('click', function() {
            sortTable(this.dataset.table, parseInt(this.dataset.col), this.dataset.sort, this);
        });
    });
});

var sortStates = {};

function sortTable(tableIdx, colIdx, sortType, header) {
    var table = document.getElementById('table-' + tableIdx);
    if (!table) return;
    var tbody = table.querySelector('tbody');
    var rows = Array.from(tbody.querySelectorAll('tr'));
    var headers = table.querySelectorAll('th');

    var key = tableIdx + '-' + colIdx;
    if (sortStates[key] === 'asc') {
        sortStates[key] = 'desc';
    } else {
        sortStates[key] = 'asc';
    }
    var ascending = sortStates[key] === 'asc';

    headers.forEach(function(h) {
        h.classList.remove('sort-asc', 'sort-desc');
        h.classList.add('sort-indicator');
    });
    header.classList.remove('sort-indicator');
    header.classList.add(ascending ? 'sort-asc' : 'sort-desc');

    rows.sort(function(a, b) {
        var valA, valB;
        if (sortType === 'text') {
            valA = a.cells[colIdx].textContent.trim().toLowerCase();
            valB = b.cells[colIdx].textContent.trim().toLowerCase();
            return ascending ? valA.localeCompare(valB) : valB.localeCompare(valA);
        } else {
            valA = parseInt(a.cells[colIdx].getAttribute('data-pass-percent') || '0', 10);
            valB = parseInt(b.cells[colIdx].getAttribute('data-pass-percent') || '0', 10);
            return ascending ? valA - valB : valB - valA;
        }
    });

    rows.forEach(function(row) { tbody.appendChild(row); });
}
</script>
</body>
</html>
`
