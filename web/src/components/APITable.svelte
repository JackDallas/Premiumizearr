<script>
  import { DataTable, InlineLoading } from "carbon-components-svelte";
  import { CalculateAPIPath } from "../Utilities/web_root";

  export let totalName = "";
  export let headers = {};
  export let updateTimeSeconds = 10;
  export let APIpath = "/api/transfers";
  export let dataToRows = function (data) {
    if (!data) return [];
    return data;
  };
  
  let updating = false;
  let status = "";
  let rows = [];
  $: statusIndicator = updating ? "active" : "finished";

  function UpdateFromAPI() {
    if (updating) return;
    // Refresh from endpoint
    updating = true;
    fetch(CalculateAPIPath(APIpath))
      .then((res) => res.json())
      .then((data) => {
        rows = dataToRows(data.data);
        status = data.status;
        updating = false;
      })
      .catch((err) => {
        console.error(err);
        updating = false;
      });
  }
  UpdateFromAPI();
  setInterval(() => {
    UpdateFromAPI();
  }, updateTimeSeconds * 1000);

  function safeLength(obj) {
    return obj ? Object.keys(obj).length : 0;
  }
</script>

<main>
  {#if totalName !== ""}
    <p>
      {totalName}
      {safeLength(rows)}
    </p>
  {/if}
  <p>
    <InlineLoading status={statusIndicator} description="Update status" />
  </p>
  <p>
    Message: {status}
  </p>
  <p>
    <DataTable sortable {headers} {rows} />
  </p>
</main>
