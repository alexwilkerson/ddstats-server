<template>
  <div class="wrapper">
    <h1>{{ playerName }} ({{ $route.params.id }})</h1>
    <div>
      <RecentPlayerGamesTable @onPlayerNameLoad="onPlayerNameLoad" />
    </div>
  </div>
</template>

<script>
import RecentPlayerGamesTable from "../components/RecentPlayerGamesTable";
import "vue-select/dist/vue-select.css";

export default {
  data() {
    return {
      playerName: ""
    };
  },
  components: {
    RecentPlayerGamesTable
  },
  methods: {
    onPlayerNameLoad(name) {
      this.playerName = name;
    }
  }
};
</script>

<style scoped>
.wrapper {
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 40px;
  padding-bottom: 40px;
  max-width: 800px;
  margin: auto;
}
h1 {
  text-align: center;
  margin-bottom: 20px;
  color: var(--v-primary-base);
}
.v-select {
  max-width: 400px;
  margin: 0 auto 20px auto;
}
</style>