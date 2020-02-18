<template>
  <v-card>
    <v-toolbar flat color="th" dark dense>
      <v-toolbar-title>Player Activity</v-toolbar-title>
    </v-toolbar>
    <v-tabs vertical>
      <v-tab v-if="data.devil_daggers_list.length > 0">
        <v-icon left fill="#c33409">$dagger</v-icon>
        New Devil Daggers ({{ data.since_devil_daggers.toLocaleString() }})
      </v-tab>
      <v-tab v-if="data.gold_daggers_list.length > 0">
        <v-icon left fill="#ffcd00">$dagger</v-icon>
        New Gold Daggers ({{ data.since_gold_daggers.toLocaleString() }})
      </v-tab>
      <v-tab v-if="data.silver_daggers_list.length > 0">
        <v-icon left fill="#acacac">$dagger</v-icon>
        New Silver Daggers ({{ data.since_silver_daggers.toLocaleString() }})
      </v-tab>
      <v-tab v-if="data.bronze_daggers_list.length > 0">
        <v-icon left fill="#ff8300">$dagger</v-icon>
        New Bronze Daggers ({{ data.since_bronze_daggers.toLocaleString() }})
      </v-tab>
      <v-tab
        v-if="data.active_players_list.length > 0"
      >Active Players ({{ data.active_players.toLocaleString() }})</v-tab>
      <v-tab
        v-if="data.new_players_list.length > 0"
      >New Players ({{ data.new_players.toLocaleString() }})</v-tab>

      <v-tab-item v-if="data.devil_daggers_list.length > 0">
        <v-card flat>
          <v-data-table
            :headers="headers"
            :items="data.devil_daggers_list"
            item-key="player_id"
            disableSort
            disableFiltering
            :hide-default-footer="data.devil_daggers_list.length < 11"
            :footer-props="{
              itemsPerPageOptions: [10],
              showCurrentPage: true
            }"
            dense
          ></v-data-table>
        </v-card>
      </v-tab-item>
      <v-tab-item v-if="data.gold_daggers_list.length > 0">
        <v-card flat>
          <v-data-table
            :headers="headers"
            :items="data.gold_daggers_list"
            item-key="player_id"
            disableSort
            disableFiltering
            :hide-default-footer="data.gold_daggers_list.length < 11"
            :footer-props="{
              itemsPerPageOptions: [10],
              showCurrentPage: true
            }"
            dense
          ></v-data-table>
        </v-card>
      </v-tab-item>
      <v-tab-item v-if="data.silver_daggers_list.length > 0">
        <v-card flat>
          <v-data-table
            :headers="headers"
            :items="data.silver_daggers_list"
            item-key="player_id"
            disableSort
            disableFiltering
            :hide-default-footer="data.silver_daggers_list.length < 11"
            :footer-props="{
              itemsPerPageOptions: [10],
              showCurrentPage: true
            }"
            dense
          ></v-data-table>
        </v-card>
      </v-tab-item>
      <v-tab-item v-if="data.bronze_daggers_list.length > 0">
        <v-card flat>
          <v-data-table
            :headers="headers"
            :items="data.bronze_daggers_list"
            item-key="player_id"
            disableSort
            disableFiltering
            :hide-default-footer="data.bronze_daggers_list.length < 11"
            :footer-props="{
              itemsPerPageOptions: [10],
              showCurrentPage: true
            }"
            dense
          ></v-data-table>
        </v-card>
      </v-tab-item>
      <v-tab-item v-if="data.active_players_list.length > 0">
        <v-card flat>
          <v-data-table
            :headers="headers"
            :items="data.active_players_list"
            item-key="player_id"
            disableSort
            disableFiltering
            :hide-default-footer="data.active_players_list.length < 11"
            :footer-props="{
              itemsPerPageOptions: [10],
              showCurrentPage: true
            }"
            dense
          ></v-data-table>
        </v-card>
      </v-tab-item>
      <v-tab-item v-if="data.new_players_list.length > 0">
        <v-card flat>
          <v-data-table
            :headers="headers"
            :items="data.new_players_list"
            item-key="player_id"
            disableSort
            disableFiltering
            :hide-default-footer="data.new_players_list.length < 11"
            :footer-props="{
              itemsPerPageOptions: [10],
              showCurrentPage: true
            }"
            dense
          ></v-data-table>
        </v-card>
      </v-tab-item>
    </v-tabs>
  </v-card>
</template>

<script>
export default {
  props: ["data"],
  data() {
    return {
      accuracyHover: false,
      headers: [
        { text: "Rank", value: "rank" },
        { text: "Name", value: "player_name" },
        { text: "Game Time", value: "game_time" }
      ],
      options: {
        page: 1,
        rowsPerPage: 10
      }
    };
  }
};
</script>

<style>
.since-text {
  color: #aeaeae;
  font-family: "alte_haas_grotesk", "Helvetica Neue", Helvetica, Arial;
  font-style: oblique;
  font-size: 12px !important;
}
</style>