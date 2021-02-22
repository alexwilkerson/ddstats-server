<template>
  <div class="wrapper-dapper">
    <h1 :style="{marginBottom: '20px'}">
      <v-icon
        v-if="$root.checkPlayerLive($route.params.id)"
        class="icon online-green"
      >mdi-access-point</v-icon>
      Currently: {{ $root.status }} {{ $root.status === "Dead" ? `(${getDeathType($root.state.death_type)})` : "" }}
    </h1>
    <v-card :style="{padding: '12px'}">
      <v-simple-table dense>
        <template v-slot:default>
          <tbody>
            <tr>
              <td class="text-left">Game Time</td>
              <td
                class="text-right"
                v-if="$root.state.game_time !== undefined"
              >{{ $root.state.game_time.toFixed(4) }}s</td>
              <td class="text-right" v-else>0000.0s</td>
            </tr>
            <tr>
              <td class="text-left">Gems</td>
              <td class="text-right" v-if="$root.state.gems !== undefined">{{ $root.state.gems }}</td>
              <td class="text-right" v-else>0</td>
            </tr>
            <tr>
              <td class="text-left">Homing Daggers</td>
              <td
                class="text-right"
                v-if="$root.state.homing_daggers !== undefined"
              >{{ $root.state.homing_daggers }}</td>
              <td class="text-right" v-else>0</td>
            </tr>
            <tr>
              <td class="text-left">Enemies Alive</td>
              <td
                class="text-right"
                v-if="$root.state.enemies_alive !== undefined"
              >{{ $root.state.enemies_alive }}</td>
              <td class="text-right" v-else>0</td>
            </tr>
            <tr>
              <td class="text-left">Enemies Killed</td>
              <td
                class="text-right"
                v-if="$root.state.enemies_killed !== undefined"
              >{{ $root.state.enemies_killed }}</td>
              <td class="text-right" v-else>0</td>
            </tr>
            <tr>
              <td class="text-left">Accuracy</td>
              <td class="text-right" v-if="$root.state.daggers_fired !== undefined">
                <span v-if="$root.state.daggers_fired > 0">
                  {{
                  (($root.state.daggers_hit / $root.state.daggers_fired) * 100).toFixed(
                  2
                  )
                  }}%
                </span>
                <span v-else>0.00%</span>
              </td>
              <td class="text-right" v-else>0.00%</td>
            </tr>
            <tr v-if="$root.state.level_two_time !== undefined && $root.state.level_two_time != 0">
              <td class="text-left">Level 2</td>
              <td class="text-right">{{ $root.state.level_two_time.toFixed(4) }}s</td>
            </tr>
            <tr
              v-if="$root.state.level_three_time !== undefined && $root.state.level_three_time != 0"
            >
              <td class="text-left">Level 3</td>
              <td class="text-right">{{ $root.state.level_three_time.toFixed(4) }}s</td>
            </tr>
            <tr
              v-if="$root.state.level_four_time !== undefined && $root.state.level_four_time != 0"
            >
              <td class="text-left">Level 4</td>
              <td class="text-right">{{ $root.state.level_four_time.toFixed(4) }}s</td>
            </tr>
            <tr v-if="$root.state.levi_down_time !== undefined && $root.state.levi_down_time != 0">
              <td class="text-left">Levi Down Time</td>
              <td class="text-right">{{ $root.state.levi_down_time.toFixed(4) }}s</td>
            </tr>
            <tr v-if="$root.state.orb_down_time !== undefined && $root.state.orb_down_time != 0">
              <td class="text-left">Orb Down Time</td>
              <td class="text-right">{{ $root.state.orb_down_time.toFixed(4) }}s</td>
            </tr>
            <tr>
              <td class="text-left">Viewers</td>
              <td class="text-right">{{ $root.watchers }}</td>
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
  </div>
</template>

<script>
export default {
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
  "fallen",
  "swarmed",
  "impaled",
  "gored",
  "infested",
  "opened",
  "purged",
  "desecrated",
  "sacrificed",
  "eviscerated",
  "annihilated",
  "intoxicated",
  "envenmonated",
  "incarnated",
  "discarnated",
  "barbed"
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
