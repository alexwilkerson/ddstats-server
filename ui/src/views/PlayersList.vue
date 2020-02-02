<template>
  <div class="wrapper">
    <h1>Players</h1>
    <div>
      <PlayersTable />
    </div>
  </div>
</template>

<script>
import PlayersTable from "../components/PlayersTable";

export default {
  components: {
    PlayersTable
  }
};
</script>

<style scoped>
.wrapper {
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 40px;
  padding-bottom: 40px;
  max-width: 860px;
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
