<template>
  <div>
    <v-subheader>Provider types</v-subheader>
    <div class="elevation-1" id="admin-provider">
      <v-toolbar flat dense color="white">
        <v-spacer></v-spacer>

        <v-dialog v-model="dialog" max-width="600px" :scrollable="true">
          <v-btn small slot="activator" color="primary" dark>New</v-btn>
          <v-card>
            <v-card-title>
              <span class="headline">{{ this.index === -1 ? "New provider" : "Edit provider" }}</span>
            </v-card-title>

            <v-card-text max-width="600px">
              <v-alert :value="dialogError" color="error" icon="warning" outline v-html="dialogError"></v-alert>

              <v-container grid-list-sm id="dialog">
                <v-layout wrap>
                  <v-flex xs12><v-divider></v-divider><v-subheader>General settings</v-subheader></v-flex>
                  <v-flex sm12 md6>
                    <v-text-field
                      v-model="editing.name"
                      prepend-icon="create"
                      label="Provider name">
                    </v-text-field>
                  </v-flex>
                  <v-flex sm12 md6>
                    <v-text-field
                      v-model="editing.key"
                      prepend-icon="create"
                      label="Provider key">
                    </v-text-field>
                  </v-flex>
                  <v-flex xs12>
                    <v-text-field
                      v-model="editing.file_link_format"
                      prepend-icon="code"
                      label="Provider file link format">
                    </v-text-field>
                  </v-flex>
                  <v-flex xs12>
                    <v-text-field
                      v-model="editing.uri"
                      prepend-icon="open_in_browser"
                      label="Provider uri">
                    </v-text-field>
                  </v-flex>
                  <v-flex xs12>
                    <v-textarea
                      v-model="editing.description"
                      prepend-icon="short_text"
                      rows="1"
                      auto-grow
                      label="Provider description">
                    </v-textarea>
                  </v-flex>
                  <v-flex xs12>
                    <v-checkbox
                      v-model="editing.case_required"
                      prepend-icon="format_size"
                      label="Case sensitive URI"
                      messages="Some VCS provider services make usernames and provider names case sensitive.">
                    </v-checkbox>
                  </v-flex>

                  <v-flex xs12><v-divider class="mt-4"></v-divider><v-subheader>Project settings</v-subheader></v-flex>
                  <v-flex xs12>
                    <v-select
                      v-model="editing.outdated_duration_hours"
                      prepend-icon="timelapse"
                      :items="outdated_duration_types"
                      item-text="name"
                      item-value="value"
                      required
                      label="Timespan at which a project is considered &quot;outdated&quot;"
                      messages="If automatic rescanning is enabled, this will trigger rescanning after this period of time. If automatic rescanning is disabled, this will cause rescanning if a user visits a repository that is this old.">
                    </v-select>
                  </v-flex>
                  <v-flex xs12>
                    <v-checkbox
                      v-model="editing.automatic_rescan"
                      prepend-icon="refresh"
                      label="Automatic rescanning of outdated projects"
                      messages="When enabled, a scheduler will periodically check for any repositories that haven't been updated in awhile (and if so, will rescan them). If disabled, when a users loads the page, it will trigger a rescan if it's outdated. Users will still be able to manually scan.">
                    </v-checkbox>
                  </v-flex>
                  <v-flex xs12>
                    <v-checkbox
                      v-model="editing.recursive_fetch"
                      prepend-icon="format_indent_increase"
                      label="Recursive fetching when scanning"
                      messages="E.g. git submodules.">
                    </v-checkbox>
                  </v-flex>

                  <v-flex xs12><v-divider class="mt-4"></v-divider><v-subheader>Abuse &amp; rate-limiting settings</v-subheader></v-flex>
                  <v-flex xs12>
                    <v-select
                      v-model="editing.hourly_ratelimit"
                      prepend-icon="show_chart"
                      :items="hourly_ratelimit_types"
                      item-text="name"
                      item-value="value"
                      required
                      label="Maximum number of scans per hour"
                      messages="Limit the total number of scans done toward repositories in this provider hourly (both manual and automatic scans would be limited by this).">
                    </v-select>
                  </v-flex>
                  <v-flex xs12>
                    <v-checkbox
                      v-model="editing.auth_required"
                      prepend-icon="security"
                      label="Must be logged in to scan/rescan projects"
                      messages="Can potentially cut down on abusive behavior by slowing down scan attempts.">
                    </v-checkbox>
                  </v-flex>
                </v-layout>
              </v-container>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" flat @click="close">Cancel</v-btn>
              <v-btn color="blue darken-1" flat @click="save">Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-toolbar>

      <v-divider></v-divider>

      <v-data-table :headers="headers" :loading="loading" :items="items">
        <template slot="items" slot-scope="props">
          <td>{{ props.item.name }}</td>
          <td>{{ props.item.key }}</td>
          <td>{{ props.item.description }}</td>
          <td><code>{{ props.item.file_link_format }}</code></td>
          <td><code>{{ props.item.uri }}</code></td>
          <td>{{ props.item.auth_required ? "Yes" : "No" }}</td>

          <td class="justify-center layout px-0">
            <v-icon small class="mr-2" @click="editItem(props.item)">edit</v-icon>
            <v-icon small @click="deleteItem(props.item)">delete</v-icon>
          </td>
        </template>

        <template slot="no-data" v-if="!loading">
            <v-alert v-if="!error" color="info" icon="info" outline :value="true">There are no providers setup currently.</v-alert>
            <v-alert v-if="error" color="error" icon="warning" outline :value="true">{{ error }}</v-alert>
        </template>
      </v-data-table>
    </div>
  </div>
</template>

<script>
export default {
  name: "admin-providers",
  data: () => ({
    dialog: false,
    loading: true,
    // TODO: show dialog for the form errors.
    error: null,
    dialogError: null,
    headers: [
      { text: "Provider name", align: "left", value: "name" },
      { text: "Key", value: "key" },
      { text: "Description", value: "description" },
      { text: "File link format", value: "file_link_format" },
      { text: "URI", value: "uri" },
      { text: "Auth required", value: "auth_required" },
      { text: "Actions", sortable: false }
    ],
    items: [],
    index: -1,
    editing: {},
    template: {
      id: "",
      name: "",
      key: "",
      description: "",
      file_link_format: "",
      uri: "",
      case_required: false,
      automatic_rescan: true,
      outdated_duration_hours: 72,
      recursive_fetch: true,
      auth_required: false,
      hourly_ratelimit: 100
    },
    outdated_duration_types: [
      { value: 1, name: "1 hour" },
      { value: 6, name: "6 hours" },
      { value: 12, name: "12 hours" },
      { value: 24, name: "1 day" },
      { value: 72, name: "3 days" },
      { value: 168, name: "1 week" },
      { value: 336, name: "2 weeks" },
      { value: 720, name: "30 days" }
    ],
    hourly_ratelimit_types: [10,50,100,250,500,1000,1500,2000,3000,5000]
  }),

  watch: {
    dialog: function(val) {
      if (typeof this.editing.id == 'undefined') {
        this.editing = Object.assign({}, this.template);
      }
      val || this.close();
    }
  },

  created: function() {
    this.initialize();
  },

  methods: {
    initialize: function() {
      this.editing = Object.assign({}, this.template);

      this.$http.get("/api/v1/admin/providers").then((resp) => {
        this.items = resp.data.providers;
        this.provider_types = resp.data.provider_types;
      }).catch((error) => {
        // TODO: mixin to handle errors?
        this.error = error.response.data.error || error;
        if (error.response.data.type) {
            this.error = error.response.data.type + ": " + this.error
        }
      }).finally(() => {
        this.loading = false
      })
    },

    editItem: function(item) {
      this.index = this.items.indexOf(item);
      this.editing = Object.assign({}, item);
      this.dialog = true;
    },

    deleteItem: function(item) {
      const index = this.items.indexOf(item);
      if (confirm("Are you sure you want to delete this repository?")) {
        this.$http.delete("/api/v1/admin/providers/" + this.items[index].id).then((resp) => {
          this.items.splice(index, 1);

          this.$snackbar.show("Successfully deleted " + item.name);
        }).catch((error) => {
          this.$snackbar.show(error.response.data.error ? error.response.data.error.replace(/\n/g, "<br/>") : error);
        })
      }
    },

    close: function() {
      this.dialog = false;
      setTimeout(() => {
        this.editing = Object.assign({}, this.template);
        this.index = -1;
      }, 300);
    },

    save: function() {
      this.dialogError = null;

      if (this.index > -1) {
        this.$http.put("/api/v1/admin/providers/" + this.editing.id, this.editing).then((resp) => {
          Object.assign(this.items[this.index], resp.data);

          this.$snackbar.show("Successfully updated " + this.editing.name);
          this.close();
        }).catch((error) => {
            this.dialogError = error.response.data.error ? error.response.data.error.replace(/\n/g, "<br/>") : error;
          if (error.response.data.type) {
            this.dialogError = error.response.data.type + ": " + this.dialogError
          }
        })
      } else {
        this.$http.post("/api/v1/admin/providers", this.editing).then((resp) => {
          this.items.push(resp.data);

          this.$snackbar.show("Successfully added repository");
          this.close();
        }).catch((error) => {
            this.dialogError = error.response.data.error ? error.response.data.error.replace(/\n/g, "<br/>") : error;
          if (error.response.data.type) {
            this.dialogError = error.response.data.type + ": " + this.dialogError
          }
        })
      }
    }
  }
};
</script>

<style scoped>
#admin-repository .v-toolbar__content { padding: 5px; }
#dialog { padding: 0; }
.v-card__title { padding: 16px 16px 0 16px; }
</style>
