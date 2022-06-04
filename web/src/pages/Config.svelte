<script>
  import {
    Row,
    Column,
    Button,
    TextInput,
    Modal,
    FormGroup,
    Dropdown,
  } from "carbon-components-svelte";
  import {
    Save,
    CheckmarkFilled,
    AddFilled,
    TrashCan,
  } from "carbon-icons-svelte";

  let webRoot = window.location.href;

  let config = {
    BlackholeDirectory: "",
    DownloadsDirectory: "",
    UnzipDirectory: "",
    BindIP: "",
    BindPort: "",
    WebRoot: "",
    SimultaneousDownloads: 0,
    Arrs: [],
  };

  let inputDisabled = true;

  let errorModal = false;
  let errorMessage = "";

  let saveIcon = Save;

  function getConfig() {
    inputDisabled = true;
    fetch(webRoot + "api/config")
      .then((response) => response.json())
      .then((data) => {
        config = data;
        inputDisabled = false;
      })
      .catch((error) => {
        console.error("Error: ", error);
      });
  }

  function submit() {
    inputDisabled = true;
    fetch(webRoot + "api/config", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(config),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.succeeded) {
          saveIcon = CheckmarkFilled;
          getConfig();
          setTimeout(() => {
            saveIcon = Save;
          }, 1000);
        } else {
          errorMessage = data.status;
          errorModal = true;
          getConfig();
        }
      })
      .catch((error) => {
        console.error("Error: ", error);
        errorModal = true;
        errorMessage = error;
        setTimeout(() => {
          getConfig();
        }, 1500);
      });
  }

  function AddArr() {
    config.Arrs.push({
      Name: "New Arr",
      URL: "http://localhost:1234",
      APIKey: "xxxxxxxx",
      Type: "Sonarr",
    });
    //Force re-paint
    config.Arrs = [...config.Arrs];
  }

  function RemoveArr(index) {
    console.log(index);
    config.Arrs.splice(index, 1);
    //Force re-paint
    config.Arrs = [...config.Arrs];
  }

  getConfig();
</script>

<main>
  <Row>
    <Column>
      <h4>*Arr Settings</h4>
      <FormGroup>
        {#if config.Arrs !== undefined}
          {#each config.Arrs as arr, i}
            <h5>- {arr.Name ? arr.Name : i}</h5>
            <FormGroup>
              <TextInput
                labelText="Name"
                bind:value={arr.Name}
                disabled={inputDisabled}
              />
              <TextInput
                labelText="URL"
                bind:value={arr.URL}
                disabled={inputDisabled}
              />
              <TextInput
                labelText="APIKey"
                bind:value={arr.APIKey}
                disabled={inputDisabled}
              />
              <Dropdown
                titleText="Type"
                selectedId={arr.Type}
                on:select={(e) => {
                  config.Arrs[i].Type = e.detail.selectedId;
                }}
                items={[
                  { id: "Sonarr", text: "Sonarr" },
                  { id: "Radarr", text: "Radarr" },
                ]}
                disabled={inputDisabled}
              />
              <Button
                style="margin-top: 10px;"
                on:click={() => {
                  RemoveArr(i);
                }}
                kind="danger"
                icon={TrashCan}
                iconDescription="Delete Arr"
              />
            </FormGroup>
          {/each}
        {/if}
      </FormGroup>
      <Button on:click={AddArr} disabled={inputDisabled} icon={AddFilled}>
        Add Arr
      </Button>
    </Column>
    <Column>
      <h4>Directory Settings</h4>
      <FormGroup>
        <TextInput
          disabled={inputDisabled}
          labelText="Blackhole Directory"
          bind:value={config.BlackholeDirectory}
        />
        <TextInput
          disabled={inputDisabled}
          labelText="Download Directory"
          bind:value={config.DownloadsDirectory}
        />
        <TextInput
          disabled={inputDisabled}
          labelText="Unzip Directory"
          bind:value={config.UnzipDirectory}
        />
      </FormGroup>
      <h4>Web Server Settings</h4>
      <FormGroup>
        <TextInput
          disabled={inputDisabled}
          labelText="Bind IP"
          bind:value={config.BindIP}
        />
        <TextInput
          disabled={inputDisabled}
          labelText="Bind Port"
          bind:value={config.BindPort}
        />
        <TextInput
          disabled={inputDisabled}
          labelText="Web Root"
          bind:value={config.WebRoot}
        />
      </FormGroup>
      <h4>Download Settings</h4>
      <FormGroup>
        <TextInput
          type="number"
          disabled={inputDisabled}
          labelText="Simultaneous Downloads"
          bind:value={config.SimultaneousDownloads}
        />
      </FormGroup>
      <Button on:click={submit} icon={saveIcon} disabled={inputDisabled}
        >Save</Button
      >
    </Column>
  </Row>
</main>

<Modal
  bind:open={errorModal}
  on:open={errorModal}
  passiveModal
  modalHeading="Error Saving Config"
  on:close={() => {
    errorModal = false;
  }}
>
  <p>{errorMessage}</p>
</Modal>
