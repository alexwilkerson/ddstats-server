<template>
  <v-card>
    <v-card-title class="dd-card-title" v-if="$root.players.length === 0">No Live Players</v-card-title>
    <v-card-title class="dd-card-title dd-card-title-bottom-border" v-else>Live Players</v-card-title>
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