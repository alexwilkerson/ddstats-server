<template>
  <div class="wrapper">
    <h1>Recently Played Games</h1>
    <div>
      <RecentGamesTable />
    </div>
  </div>
</template>

<script>
import RecentGamesTable from "../components/RecentGamesTable";
import "vue-select/dist/vue-select.css";

export default {
  components: {
    RecentGamesTable
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