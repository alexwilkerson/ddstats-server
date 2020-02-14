<template>
  <div class="wrapper-dapper">
    <h1
      :style="{marginBottom: '20px'}"
    >Currently: {{ $root.status }} {{ $root.status === "Dead" ? `(${getDeathType($root.state.death_type)})` : "" }}</h1>
    <v-simple-table dense>
      <template v-slot:default>
        <tbody>
          <tr>
            <td class="text-left">Game Time</td>
            <td class="text-right">{{ $root.state.game_time.toFixed(4) }}s</td>
          </tr>
          <tr>
            <td class="text-left">Gems</td>
            <td class="text-right">{{ $root.state.gems }}</td>
          </tr>
          <tr>
            <td class="text-left">Homing Daggers</td>
            <td class="text-right">{{ $root.state.homing_daggers }}</td>
          </tr>
          <tr>
            <td class="text-left">Enemies Alive</td>
            <td class="text-right">{{ $root.state.enemies_alive }}</td>
          </tr>
          <tr>
            <td class="text-left">Enemies Killed</td>
            <td class="text-right">{{ $root.state.enemies_killed }}</td>
          </tr>
          <tr>
            <td class="text-left">Accuracy</td>
            <td class="text-right">
              <span v-if="$root.state.daggers_fired > 0">
                {{
                (($root.state.daggers_hit / $root.state.daggers_fired) * 100).toFixed(
                2
                )
                }}%
              </span>
              <span v-else>0%</span>
            </td>
          </tr>
          <tr v-if="$root.state.level_two_time != 0">
            <td class="text-left">Level 2</td>
            <td class="text-right">{{ $root.state.level_two_time.toFixed(4) }}s</td>
          </tr>
          <tr v-if="$root.state.level_three_time != 0">
            <td class="text-left">Level 3</td>
            <td class="text-right">{{ $root.state.level_three_time.toFixed(4) }}s</td>
          </tr>
          <tr v-if="$root.state.level_four_time != 0">
            <td class="text-left">Level 4</td>
            <td class="text-right">{{ $root.state.level_four_time.toFixed(4) }}s</td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </div>
</template>

<script>
export default {
  data() {
    return {
      state: {}
    };
  },
  methods: {
    getDeathType(id) {
      if (id === -1) {
        return "N/A";
      }
      return deathTypes[id];
    }
  }
};
const deathTypes = [
  "FALLEN",
  "SWARMED",
  "IMPALED",
  "GORED",
  "INFESTED",
  "OPENED",
  "PURGED",
  "DESECRATED",
  "SACRIFICED",
  "EVISCERATED",
  "ANNIHILATED",
  "INTOXICATED",
  "ENVENMONATED",
  "INCARNATED",
  "DISCARNATED",
  "BARBED"
];
</script>

<style>
.wrapper-dapper {
  margin: 0 auto 20px auto;
  max-width: 650px !important;
  text-align: center;
  font-family: "alte_haas_grotesk_bold", "Helvetica Neue", Helvetica, Arial;
}
.wrapper-dapper h1 {
  color: var(--v-primary-base);
}
</style>
