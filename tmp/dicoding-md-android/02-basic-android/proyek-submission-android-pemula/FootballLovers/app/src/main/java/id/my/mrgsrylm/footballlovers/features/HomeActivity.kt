package id.my.mrgsrylm.footballlovers.features

import android.content.Intent
import android.os.Bundle
import android.view.Menu
import android.view.MenuItem
import androidx.appcompat.app.AppCompatActivity
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import id.my.mrgsrylm.footballlovers.R
import id.my.mrgsrylm.footballlovers.data.dummy.teamsDummy
import id.my.mrgsrylm.footballlovers.data.model.Team
import id.my.mrgsrylm.footballlovers.features.about.AboutActivity
import id.my.mrgsrylm.footballlovers.features.team.TeamAdapter

class HomeActivity : AppCompatActivity() {
    private lateinit var rvTeam: RecyclerView
    private var teams: ArrayList<Team> = arrayListOf()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_home)

        val actionBar = supportActionBar
        actionBar!!.title = getString(R.string.home)

        rvTeam = findViewById(R.id.rv_team_tile)
        rvTeam.setHasFixedSize(true)

        teams.addAll(teamsDummy)
        showTeamsTile()
    }

    override fun onCreateOptionsMenu(menu: Menu?): Boolean {
        menuInflater.inflate(R.menu.menus, menu)
        return super.onCreateOptionsMenu(menu)
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        setMode(item.itemId)
        return super.onOptionsItemSelected(item)
    }

    private fun setMode(selectedMode: Int) {
        when (selectedMode) {
            R.id.miCompose -> {
                val iAbout = Intent(
                    this@HomeActivity,
                    AboutActivity::class.java
                )
                startActivity(iAbout)
            }
        }
    }

    private fun showTeamsTile() {
        rvTeam.layoutManager = LinearLayoutManager(this)
        val adapter = TeamAdapter(teams)
        rvTeam.adapter = adapter
    }
}