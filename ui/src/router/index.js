import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: Home
  },
  {
    path: "/games/:id",
    name: "games",
    props: route => {
      const id = Number.parseInt(route.params.id, 10);
      if (Number.isNaN(id)) {
        return 0;
      }
      return { id };
    },
    component: () =>
      import(/* webpackChunkName: "games" */ "../views/Games.vue")
  },
  {
    path: "/games",
    name: "gamesList",
    component: () =>
      import(/* webpackChunkName: "gamesList" */ "../views/GamesList.vue")
  },
  {
    path: "/leaderboard/:name",
    name: "leaderboard",
    component: () =>
      import(/* webpackChunkName: "leaderboard" */ "../views/Leaderboard.vue")
  },
  {
    path: "/leaderboard",
    name: "leaderboard",
    component: () =>
      import(/* webpackChunkName: "leaderboard" */ "../views/Leaderboard.vue")
  },
  // {
  //   path: "/leaderboard/:name",
  //   name: "leaderboardByName",
  //   component: () =>
  //     import(
  //       /* webpackChunkName: "leaderboardByName" */ "../views/LeaderboardByName.vue"
  //     )
  // },
  {
    path: "/players",
    name: "playersList",
    component: () =>
      import(/* webpackChunkName: "playersList" */ "../views/PlayersList.vue")
  },
  {
    path: "/players/:id",
    name: "players",
    props: route => {
      const id = Number.parseInt(route.params.id, 10);
      if (Number.isNaN(id)) {
        return 0;
      }
      return { id };
    },
    component: () =>
      import(/* webpackChunkName: "players" */ "../views/Players.vue")
  },
  {
    path: "/info",
    name: "info",
    component: () => import(/* webpackChunkName: "info" */ "../views/Info.vue")
  },
  {
    path: "/download",
    name: "download",
    component: () =>
      import(/* webpackChunkName: "download" */ "../views/Download.vue")
  },
  {
    path: "**",
    name: "notFound",
    component: () =>
      import(/* webpackChunkName: "notFound" */ "../views/NotFound.vue")
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
