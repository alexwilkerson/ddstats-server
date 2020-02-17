<template>
  <v-card>
    <v-card-title class="dd-card-title dd-card-title-bottom-border">News</v-card-title>
    <p class="dd-card-list-item"></p>
    <p
      v-if="data.devil_daggers_list.length > 0"
      class="dd-card-list-item"
    >Congratulations to {{ newDevilDaggerers() }} on passing 500 seconds and getting their devil daggers!
    <ul>
      <li v-for="player in data.devil_daggers_list" :key="player.player_id"><strong>{{ player.player_name }}</strong> - {{ player.game_time }}s</li>
    </ul>
    </p>

  </v-card>
</template>

<script>
export default {
  props: ["data"],
  methods: {
    newDevilDaggerers() {
      let users = [];
      for (let i = 0; i < this.data.devil_daggers_list.length; i++) {
        users.push(this.data.devil_daggers_list[i].player_name);
      }
      return users.join(", ").replace(/,(?!.*,)/gim, ", and");
    }
  }
};
</script>

<style>
</style>