package id.my.mrgsrylm.githubuser.ui.activity

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.ViewModelProvider
import androidx.recyclerview.widget.LinearLayoutManager
import id.my.mrgsrylm.githubuser.R
import id.my.mrgsrylm.githubuser.data.local.SettingPreferences
import id.my.mrgsrylm.githubuser.data.local.dataStore
import id.my.mrgsrylm.githubuser.data.remote.response.UserItem
import id.my.mrgsrylm.githubuser.databinding.ActivityMainBinding
import id.my.mrgsrylm.githubuser.helper.PreferencesViewModelFactory
import id.my.mrgsrylm.githubuser.ui.adapter.ItemUserAdapter
import id.my.mrgsrylm.githubuser.ui.viewmodel.MainViewModel

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var adapter: ItemUserAdapter

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val pref = SettingPreferences.getInstance(application.dataStore)
        val mainViewModel = ViewModelProvider(this, PreferencesViewModelFactory(pref)).get(
            MainViewModel::class.java
        )

        mainViewModel.users.observe(this) { users ->
            setUsersData(users)
        }

        adapter = ItemUserAdapter()
        binding.rvUser.adapter = adapter
        val layoutManager = LinearLayoutManager(this)
        binding.rvUser.layoutManager = layoutManager

        mainViewModel.isLoading.observe(this) {
            showLoading(it)
        }

        with(binding) {
            searchView.setupWithSearchBar(searchBar)
            searchView
                .editText
                .setOnEditorActionListener { textView, actionId, event ->
                    val query = searchView.text.toString()
                    searchBar.setText(searchView.text)
                    searchView.hide()
                    mainViewModel.searchUser(query)
                    false
                }
            searchBar.inflateMenu(R.menu.option_menu)
            searchBar.setOnMenuItemClickListener { menuItem ->
                when (menuItem.itemId) {
                    R.id.menu_favorite -> {
                        val intent = Intent(this@MainActivity, FavoriteActivity::class.java)
                        startActivity(intent)
                        true
                    }

                    R.id.menu_setting -> {
                        val intent = Intent(this@MainActivity, SettingActivity::class.java)
                        startActivity(intent)
                        true
                    }

                    else -> false
                }
            }
        }
    }

    private fun setUsersData(users: List<UserItem>) {
        if (users.isEmpty()) {
            binding.rvUser.visibility = View.GONE
            binding.tvNotFound.visibility = View.VISIBLE
        } else {
            binding.tvNotFound.visibility = View.GONE
            adapter.submitList(users)
        }
    }

    private fun showLoading(isLoading: Boolean) {
        if (isLoading) {
            binding.rvUser.visibility = View.GONE
            binding.tvNotFound.visibility = View.GONE
            binding.progressBar.visibility = View.VISIBLE
        } else {
            binding.progressBar.visibility = View.GONE
            binding.rvUser.visibility = View.VISIBLE
            binding.tvNotFound.visibility = View.GONE
        }
    }
}