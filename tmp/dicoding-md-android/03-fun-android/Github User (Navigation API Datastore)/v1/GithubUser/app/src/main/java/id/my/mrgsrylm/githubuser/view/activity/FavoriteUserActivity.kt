package id.my.mrgsrylm.githubuser.view.activity

import android.os.Bundle
import android.view.MenuItem
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.recyclerview.widget.LinearLayoutManager
import id.my.mrgsrylm.githubuser.databinding.ActivityFavoriteUserBinding
import id.my.mrgsrylm.githubuser.view.adapter.FavoriteUserAdapter
import id.my.mrgsrylm.githubuser.viewmodel.ViewModelFactory
import id.my.mrgsrylm.githubuser.viewmodel.user.FavoriteUserViewModel

class FavoriteUserActivity : AppCompatActivity() {

    private lateinit var binding: ActivityFavoriteUserBinding
    private lateinit var adapter: FavoriteUserAdapter
    private val favoriteUserViewModel by viewModels<FavoriteUserViewModel> {
        ViewModelFactory.getInstance(application)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityFavoriteUserBinding.inflate(layoutInflater)
        setContentView(binding.root)

        supportActionBar?.setDisplayHomeAsUpEnabled(true)

        favoriteUserViewModel.getUser().observe(this) { listFavoriteUser ->
            if (listFavoriteUser != null) {
                adapter.setListFavorites(listFavoriteUser)
            }
        }

        adapter = FavoriteUserAdapter()
        binding.rvFavoriteUser.layoutManager = LinearLayoutManager(this)
        binding.rvFavoriteUser.adapter = adapter

    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        when (item.itemId) {
            android.R.id.home -> {
                finish()
                return true
            }
        }
        return super.onContextItemSelected(item)
    }
}