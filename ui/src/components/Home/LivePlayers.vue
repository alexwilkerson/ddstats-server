<template>
  <v-card>
    <v-toolbar flat color="th" dark dense v-if="$root.players.length === 0">
      <v-toolbar-title>No Live Players</v-toolbar-title>
    </v-toolbar>
    <v-toolbar flat color="th" dark dense v-else>
      <v-toolbar-title>Live Players</v-toolbar-title>
    </v-toolbar>
    <v-simple-table>
      <tbody>
        <tr
          v-for="player in $root.players"
          :key="player.player_id"
          @click="selectPlayer(player)"
          class="pointer"
        >
          <td class="dd-card-list-item-bold">
            <v-icon class="icon online-green" small>mdi-access-point</v-icon>
            {{ player.player_name }}
          </td>
        </tr>
      </tbody>
    </v-simple-table>
  </v-card>
</template>

<script>
export default {
  data() {
    return {};
  },
  methods: {
    selectPlayer: function(player) {
      this.$router.push("/players/" + player.player_id);
    }
  }
};
</script>

<style>
</style>