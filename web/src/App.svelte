<script>
  import "carbon-components-svelte/css/g100.css";

  import {
    Grid,
    Row,
    Column,
    Tabs,
    Tab,
    TabContent,
  } from "carbon-components-svelte";
  import Config from "./pages/Config.svelte";
  import Info from "./pages/Info.svelte";
</script>

<main>
  <Grid fullWidth>
    <Row>
      <Column>
        <Tabs>
          <Tab label="Info" />
          <Tab label="Config" />
          <svelte:fragment slot="content">
            <TabContent><Info /></TabContent>
            <TabContent><Config /></TabContent>
          </svelte:fragment>
        </Tabs>
      </Column>
    </Row>
  </Grid>
</main>
