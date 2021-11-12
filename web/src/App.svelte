<script>
  import { mdiArrowDown, mdiArrowUp, mdiSwapVertical } from "@mdi/js";
  import {
    DataTable,
    DataTableHead,
    DataTableRow,
    DataTableCell,
    DataTableBody,
    MaterialApp,
    Icon,
    Button,
  } from "svelte-materialify";

  let theme = "dark";
  let sortField = "progress";
  let sortDirection = true;
  let currentSortSVG = mdiArrowDown;

  let transfers = [];
  let downloads = [];

  UpdateFromAPI();

  setInterval(() => {
    UpdateFromAPI();
  }, 10 * 1000);

  function FileNameFromPath(path) {
    return path.split("/").pop();
  }

  function UpdateFromAPI() {
    // Refresh from endpoints
    // transfers
    fetch("/api/transfers")
      .then((res) => res.json())
      .then((data) => {
        transfers = data;
        SortTransfersBy(sortField);
      });

    // downloads
    fetch("/api/downloads")
      .then((res) => res.json())
      .then((data) => {
        downloads = data;
        SortTransfersBy(sortField);
      });
  }

  function SortTransfersBy(field) {
    if (field == sortField) {
      sortDirection = !sortDirection;
      currentSortSVG = sortDirection ? mdiArrowDown : mdiArrowUp;
    }

    transfers = transfers.sort((a, b) => {
      if (sortDirection) {
        if (a[field] < b[field]) {
          return -1;
        }
        if (a[field] > b[field]) {
          return 1;
        }
        return 0;
      } else {
        if (a[field] > b[field]) {
          return -1;
        }
        if (a[field] < b[field]) {
          return 1;
        }
        return 0;
      }
    });
    sortField = field;
  }
</script>

<MaterialApp {theme}>
  <h2>Premiumizearrd</h2>
  <h3>Transfers</h3>
  <DataTable>
    <DataTableHead>
      <DataTableRow>
        <DataTableCell>
          Name
          {#if sortField == "name"}
            <Button
              icon="true"
              size="1"
              on:click={SortTransfersBy.bind(this, "name")}
            >
              <Icon path={currentSortSVG} />
            </Button>
          {:else}
            <Button
              icon="true"
              size="1"
              on:click={SortTransfersBy.bind(this, "name")}
            >
              <Icon path={mdiSwapVertical} />
            </Button>
          {/if}
        </DataTableCell>
        <DataTableCell>
          Status
          {#if sortField == "status"}
            <Button
              icon="true"
              size="1"
              on:click={SortTransfersBy.bind(this, "status")}
            >
              <Icon path={currentSortSVG} />
            </Button>
          {:else}
            <Button
              icon="true"
              size="1"
              on:click={SortTransfersBy.bind(this, "status")}
            >
              <Icon path={mdiSwapVertical} />
            </Button>
          {/if}
        </DataTableCell>
        <DataTableCell numeric>
          Progress
          {#if sortField == "progress"}
            <Button
              icon="true"
              size="1"
              on:click={SortTransfersBy.bind(this, "progress")}
            >
              <Icon path={currentSortSVG} />
            </Button>
          {:else}
            <Button
              icon="true"
              size="1"
              on:click={SortTransfersBy.bind(this, "progress")}
            >
              <Icon path={mdiSwapVertical} />
            </Button>
          {/if}
        </DataTableCell>
        <DataTableCell>Message</DataTableCell>
      </DataTableRow>
    </DataTableHead>
    <DataTableBody>
      {#each transfers as transfer}
        <DataTableRow>
          <DataTableCell>{transfer.name}</DataTableCell>
          <DataTableCell>{transfer.status}</DataTableCell>
          <DataTableCell numeric
            >{(transfer.progress * 100).toFixed(0)}%</DataTableCell
          >
          <DataTableCell>{transfer.message}</DataTableCell>
        </DataTableRow>
      {/each}
    </DataTableBody>
  </DataTable>
  <h3>Downloads Queue</h3>
  <DataTable>
    <DataTableHead>
      <DataTableRow>
        <DataTableCell>FileName</DataTableCell>
      </DataTableRow>
    </DataTableHead>
    <DataTableBody>
      {#each downloads as download}
        <DataTableRow>
          <DataTableCell>{FileNameFromPath(download)}</DataTableCell>
        </DataTableRow>
      {/each}
    </DataTableBody>
  </DataTable>
</MaterialApp>
