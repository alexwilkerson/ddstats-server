<template>
  <v-card :style="{ padding: '12px' }" class="grotesk">
    <div :style="{ maxWidth: '600px', margin: '0 auto', fontSize: '14px' }">
      <p>
        <span :style="{ color: '#c33409' }">
          <h3>About</h3>
        </span>
      </p>
      <p>
        DDSTATS records and saves data on the video game
        <a href="https://devildaggers.com/">Devil Daggers</a>. If you play Devil
        Daggers and you'd like to use DDSTATS, click the 'DOWNLOAD' link above.
        I started this website as a passion project in 2018. I maintain the
        software and website as much as I possibly can. I want to continue to
        support the Devil Daggers community and promote what I think is one of
        the best video games ever made. Always look for new releases and feel
        free to send me an email with any suggestions.
      </p>

      <p>
        <span :style="{ color: '#c33409' }">
          <h3>How to Use</h3>
        </span>
      </p>
      <p>
        The DDSTATS client (ddstats.exe) only works with the Windows Steam
        version of Devil Daggers. If you use Windows and the Steam version of
        Devil Daggers, download the client by clicking the 'DOWNLOAD' button of
        this website. Once you download the zip file, unzip it and simply start
        the client. Once the client is running, your participation is no longer
        required. Simply start Devil Daggers and ddstats.exe will automatically
        read all of the required data and submit it to ddstats.com.
      </p>
      <p>
        After your game has been submitted to the website, a link will appear at
        the bottom of the ddstats.exe window. You can click anywhere on the
        window and the link will automatically be copied to your clipboard. To
        exit the application, either click the X in the top right corner or
        press the F10 button.
      </p>
      <p>
        There's not really much more to it than that. There is a configuration
        file which is documentented with various options you may or may not want
        to change.
      </p>

      <p>
        <span :style="{ color: '#c33409' }">
          <h3>Custom Leaderboards</h3>
        </span>
      </p>
      <p>
        In order to have games show up in the custom leaderboard, it must first
        be registered in the database. Currently, this is a manual process. If
        you'd like me to add a custom leaderboard, play a game with the
        ddstats.exe client open, then send me a link to the game's recording on
        ddstats.com, the name of the spawnset, and if you'd like to define the
        custom dagger points for the devil dagger, gold dagger, silver dagger,
        and bronze dagger. Once I receive the request, I just add it to the
        database and the leaderboard will manage itself.
      </p>

      <p>
        <span :style="{ color: '#c33409' }">
          <h3>Contact</h3>
        </span>
      </p>
      <p>
        If for whatever reason you want to contact me (VHS) for suggestions,
        improvements, bug reports, or just to say hi, please feel free to email
        me at
        <a href="mailto:ddstats.com@gmail.com?Subject=Hi"
          >ddstats.com@gmail.com</a
        >. You can also get ahold of me in the various Devil Daggers related
        Discord channels under the name VHS. I should respond to you just as
        soon as I get the chance.
      </p>
    </div>
  </v-card>
</template>

<script>
export default {};
</script>

<style></style>
