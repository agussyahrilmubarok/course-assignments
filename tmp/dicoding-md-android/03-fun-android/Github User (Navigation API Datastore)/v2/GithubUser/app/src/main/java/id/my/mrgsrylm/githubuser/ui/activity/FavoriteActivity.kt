package id.my.mrgsrylm.githubuser.ui.activity

import android.os.Bundle
import android.view.View
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.recyclerview.widget.LinearLayoutManager
import id.my.mrgsrylm.githubuser.databinding.ActivityFavoriteBinding
import id.my.mrgsrylm.githubuser.helper.ViewModelFactory
import id.my.mrgsrylm.githubuser.ui.adapter.FavoriteUserAdapter
import id.my.mrgsrylm.githubuser.ui.viewmodel.FavoriteViewModel

class FavoriteActivity : AppCompatActivity() {

    private lateinit var binding: ActivityFavoriteBinding
    private lateinit var adapter: FavoriteUserAdapter
    private val favoriteViewModel by viewModels<FavoriteViewModel> {
        ViewModelFactory.getInstance(application)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityFavoriteBinding.inflate(layoutInflater)
        setContentView(binding.root)

        favoriteViewModel.getAll().observe(this) { favorites ->
            if (favorites != null) {
                binding.tvNotFound.visibility = View.GONE
                adapter.setListFavorites(favorites)

                if (favorites.isEmpty()) {
                    binding.tvNotFound.visibility = View.VISIBLE
                }

            }
        }

        adapter = FavoriteUserAdapter()
        val layoutManager = LinearLayoutManager(this)
        binding.rvFavorites.setHasFixedSize(true)
        binding.rvFavorites.layoutManager = layoutManager
        binding.rvFavorites.adapter = adapter
    }

//    private fun obtainViewModel(activity: AppCompatActivity): FavoriteViewModel {
//        val factory = ViewModelFactory.getInstance(activity.application)
//        return ViewModelProvider(activity, factory).get(FavoriteViewModel::class.java)
//    }
}