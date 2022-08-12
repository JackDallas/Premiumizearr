<script>
  import {
    Row,
    Column,
    Button,
    TextInput,
    Modal,
    FormGroup,
    Dropdown,
    Form,
    Checkbox,
  } from "carbon-components-svelte";
  import {
    Save,
    CheckmarkFilled,
    AddFilled,
    TrashCan,
    HelpFilled,
    MisuseOutline,
    WatsonHealthRotate_360,
  } from "carbon-icons-svelte";
  import { CalculateAPIPath } from "../Utilities/web_root";

  let config = {
    BlackholeDirectory: "",
    PollBlackholeDirectory: false,
    PollBlackholeIntervalMinutes: 10,
    DownloadsDirectory: "",
    UnzipDirectory: "",
    BindIP: "",
    BindPort: "",
    WebRoot: "",
    SimultaneousDownloads: 0,
    Arrs: [],
  };
  const ERR_SAVE = "Error Saving Config";
  const ERR_TEST = "Error Testing *arr client";

  let arrTesting = [];
  let arrTestIcons = [];
  let arrTestKind = [];

  let inputDisabled = true;

  let errorModal = false;
  let errorTitle = ERR_SAVE;
  let errorMessage = "";

  let saveIcon = Save;

  function getConfig() {
    inputDisabled = true;
    fetch(CalculateAPIPath("api/config"))
      .then((response) => response.json())
      .then((data) => {
        if (Array.isArray(data.Arrs)) {
          for (let i = 0; i < data.Arrs.length; i++) {
            SetTestArr(i, HelpFilled, "secondary", false);
          }
        }

        config = data;
        inputDisabled = false;
      })
      .catch((error) => {
        console.error("Error: ", error);
      });
  }

  function submit() {
    inputDisabled = true;
    fetch(CalculateAPIPath("api/config"), {
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
          errorTitle = ERR_SAVE;
          errorModal = true;
          getConfig();
        }
      })
      .catch((error) => {
        console.error("Error: ", error);
        errorTitle = ERR_SAVE;
        errorMessage = error;
        errorModal = true;
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
    config.Arrs.splice(index, 1);
    //Force re-paint
    config.Arrs = [...config.Arrs];
  }

  function TestArr(index) {
    SetTestArr(index, WatsonHealthRotate_360, "secondary", true);

    fetch(CalculateAPIPath("api/testArr"), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(config.Arrs[index]),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.succeeded) {
          SetTestArr(index, CheckmarkFilled, "primary", false);
          ResetArrTestDelayed(index, 10);
        } else {
          SetTestArr(index, MisuseOutline, "danger", false);
          ResetArrTestDelayed(index, 5);
          errorTitle = ERR_TEST;
          errorMessage = data.status;
          errorModal = true;
        }
      })
      .catch((error) => {
        console.error("Error: ", error);
        SetTestArr(index, MisuseOutline, "danger", false);
        ResetArrTestDelayed(index, 5);
        errorTitle = ERR_TEST;
        errorMessage = error;
        errorModal = true;
      });
  }

  function UntestArr(index) {
    SetTestArr(index, HelpFilled, "secondary", false);
  }

  function SetTestArr(index, icon, kind, testing) {
    arrTesting[index] = testing;
    arrTestIcons[index] = icon;
    arrTestKind[index] = kind;

    arrTesting = [...arrTesting];
    arrTestIcons = [...arrTestIcons];
    arrTestKind = [...arrTestKind];
  }

  function ResetArrTestDelayed(index, seconds) {
    setTimeout(() => {
      SetTestArr(index, HelpFilled, "secondary", false);
    }, 1000 * seconds);
  }

  getConfig();
</script>

<main>
  <Row>
    <Column>
      <h4>*Arr Settings</h4>
      <FormGroup>
        <TextInput
          type="number"
          disabled={inputDisabled}
          labelText="Arr Update History Interval (seconds)"
          bind:value={config.ArrHistoryUpdateIntervalSeconds}
        />
        {#if config.Arrs !== undefined}
          {#each config.Arrs as arr, i}
            <h5>- {arr.Name ? arr.Name : i}</h5>
            <FormGroup>
              <TextInput
                labelText="Name"
                bind:value={arr.Name}
                disabled={inputDisabled}
                on:input={() => {
                  UntestArr(i);
                }}
              />
              <TextInput
                labelText="URL"
                bind:value={arr.URL}
                disabled={inputDisabled}
                on:input={() => {
                  UntestArr(i);
                }}
              />
              <TextInput
                labelText="APIKey"
                bind:value={arr.APIKey}
                disabled={inputDisabled}
                on:input={() => {
                  UntestArr(i);
                }}
              />
              <Dropdown
                titleText="Type"
                selectedId={arr.Type}
                on:select={(e) => {
                  config.Arrs[i].Type = e.detail.selectedId;
                  UntestArr(i);
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
              <Button
                style="margin-top: 10px;"
                on:click={() => {
                  TestArr(i);
                }}
                disabled={arrTesting[i]}
                kind={arrTestKind[i]}
                icon={arrTestIcons[i]}
              >
                Test
              </Button>
            </FormGroup>
          {/each}
        {/if}
      </FormGroup>
      <Button on:click={AddArr} disabled={inputDisabled} icon={AddFilled}>
        Add Arr
      </Button>
    </Column>
    <Column>
      <h4>Premiumize.me Settings</h4>
      <FormGroup>
        <TextInput
          disabled={inputDisabled}
          labelText="API Key"
          bind:value={config.PremiumizemeAPIKey}
        />
      </FormGroup>
      <h4>Directory Settings</h4>
      <FormGroup>
        <TextInput
          disabled={inputDisabled}
          labelText="Blackhole Directory"
          bind:value={config.BlackholeDirectory}
        />
        <Checkbox
          disabled={inputDisabled}
          bind:checked={config.PollBlackholeDirectory}
          labelText="Poll Blackhole Directory"
        />
        <TextInput
          type="number"
          disabled={inputDisabled}
          labelText="Poll Blackhole Interval Minutes"
          bind:value={config.PollBlackholeIntervalMinutes}
        />
      </FormGroup>
      <FormGroup>
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
  modalHeading={errorTitle}
  on:close={() => {
    errorModal = false;
  }}
>
  <p>{errorMessage}</p>
</Modal>
<!-- 

{() => {
                  console.log(testStatus.get(i));
                  if (testStatus.get(i) == undefined)
                    return "secondary";
                  
                    if (testStatus.get(i) === 3) {
                    return "danger";
                  } else {
                    return "secondary";
                  }
                }}

-->
