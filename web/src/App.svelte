<script>
  import APITable from "./components/APITable.svelte";
  import "carbon-components-svelte/css/g100.css";
  import { Grid, Row, Column } from "carbon-components-svelte";
  import DateTime from "luxon";

  let dlSpeed = 0;

  let webRoot = window.location.href;

  function parseDLSpeedFromMessage(m) {
    if (m == "Loading..." || m == undefined) return 0;
    let speed = m.split(" ")[0];
    speed = speed.replace(",", "");
    let unit = m.split(" ")[1];
    if (Number.isNaN(speed)) {
      console.log("Speed is not a number: ", speed);
      console.log("Message: ", message);
      return 0;
    }
    if (unit === undefined || unit === null || unit == "") {
      console.log("Unit undefined in : " + m);
      return 0;
    } else {
      try {
        unit = unit.toUpperCase();
      } catch (error) {
        return 0;
      }
      unit = unit.replace("/", "");
      unit = unit.substring(0, 2);
      switch (unit) {
        case "KB":
          return speed * 1024;
        case "MB":
          return speed * 1024 * 1024;
        case "GB":
          return speed * 1024 * 1024 * 1024;
        default:
          console.log("Unknown unit: " + unit + " in message '" + m + "'");
          return 0;
      }
    }
  }

  function HumanReadableSpeed(bytes) {
    if (bytes < 1024) {
      return bytes + " B/s";
    } else if (bytes < 1024 * 1024) {
      return (bytes / 1024).toFixed(2) + " KB/s";
    } else if (bytes < 1024 * 1024 * 1024) {
      return (bytes / 1024 / 1024).toFixed(2) + " MB/s";
    } else {
      return (bytes / 1024 / 1024 / 1024).toFixed(2) + " GB/s";
    }
  }

  function dataToRows(data) {
    let rows = [];
    dlSpeed = 0;
    if (!data) return rows;

    for (let i = 0; i < data.length; i++) {
      let d = data[i];
      rows.push({
        id: d.id,
        name: d.name,
        status: d.status,
        progress: (d.progress * 100).toFixed(0) + "%",
        message: d.message,
      });

      let speed = parseDLSpeedFromMessage(d.message);
      if (!Number.isNaN(speed)) {
        dlSpeed += speed;
      } else {
        console.error("Invalid speed: " + d.message);
      }
    }
    return rows;
  }

  function downloadsToRows(downloads) {
    let rows = [];
    if (!downloads) return rows;

    for (let i = 0; i < downloads.length; i++) {
      let d = downloads[i];
      rows.push({
        Added: DateTime.fromMillis(d.added).toFormat('dd hh:mm:ss a'),
        name: d.name,
        progress: (d.progress * 100).toFixed(0) + "%",
      });
    }
  }
</script>

<main>
  <Grid fullWidth>
    <Row>
      <Column md={4} >
        <h3>Blackhole</h3>
        <APITable
          headers={[
            { key: "id", value: "Pos" },
            { key: "name", value: "Name", sort: false },
          ]}
          {webRoot}
          APIpath="api/blackhole"
          zebra={true}
          totalName="In Queue: "
        />
      </Column>
      <Column md={4} >
        <h3>Downloads</h3>
        <APITable
          headers={[
            { key: "added", value: "Added" },
            { key: "name", value: "Name" },
            { key: "progress", value: "Progress" },
            { key: "speed", value: "Speed" },
          ]}
          updateTimeSeconds={2}
          {webRoot}
          APIpath="api/downloads"
          zebra={true}
          totalName="Downloading: "
        />
      </Column>
    </Row>
    <Row>
      <Column>
        <h3>Transfers</h3>
        <p>Download Speed: {HumanReadableSpeed(dlSpeed)}</p>
        <APITable
          headers={[
            { key: "name", value: "Name" },
            { key: "status", value: "Status" },
            { key: "progress", value: "Progress" },
            { key: "message", value: "Message", sort: false },
          ]}
          {webRoot}
          APIpath="api/transfers"
          zebra={true}
          {dataToRows}
        />
      </Column>
    </Row>
  </Grid>
</main>
