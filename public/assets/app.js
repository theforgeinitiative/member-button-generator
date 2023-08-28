$(document).ready(function() {
    dt = $('#button-table').DataTable({
        ajax: {
            url: "/api/members",
            dataSrc: ""
        },
        buttons: [
            'selectAll',
            'selectNone'
        ],
        columns: [
            { data: null, defaultContent: "" },
            { data: "name" },
            { data: "barcode" },
            { data: "last_printed" },
        ],
        columnDefs: [
            {
                orderable: false,
                className: "select-checkbox",
                targets: 0
            }
        ],
        dom: 'Blfrtip',
        paging: false,
        searching: false,
        select: {
            style: "multi",
            selector: "td:first-child",
            toggleable: true,
        },
        order: [[1, "asc"]],
    });

    $("#search").click(function(){
        barcodes = $("#search-barcodes").val();
        dt.ajax.url("/api/members?barcodes="+encodeURIComponent(barcodes));
        dt.ajax.reload();
    }); 

    $("#reset").click(function(){
        dt.ajax.url("/api/members");
        dt.ajax.reload();
        $("#search-barcodes").val("");
    });

    $("#mark-complete").click(function(){
        rows = dt.rows({ selected: true }).every( function ( rowIdx, tableLoop, rowLoop ) {
            var data = this.data();
            $.post( "/api/members/" + data.id + "/complete" ).fail(function() {
                alert( "Failed to mark button as completed" );
                return;
              });            
        } ).remove().draw();
    });

    $("#print-buttons").click(function(){
        rows = dt.rows({ selected: true }).data().toArray()
        console.log(rows)
        fetch("/api/buttons/pdf", {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(rows),
        })
        .then(resp => resp.blob())
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.style.display = "none";
            a.href = url;
            // the filename you want
            a.download = "buttons.pdf";
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
        })
        .catch(() => alert("Failed to generate PDF"));
    });

    $("#print-labels").click(function(){
        rows = dt.rows({ selected: true }).data().toArray()
        console.log(rows)
        fetch("/api/labels/print", {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(rows),
        })
    });
    
});