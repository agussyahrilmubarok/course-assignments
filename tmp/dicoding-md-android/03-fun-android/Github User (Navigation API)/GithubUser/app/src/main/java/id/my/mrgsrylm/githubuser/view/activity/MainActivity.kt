package id.my.mrgsrylm.githubuser.view.activity

import android.app.SearchManager
import android.content.Context
import android.content.Intent
import android.os.Bundle
import android.view.Menu
import android.view.View
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.widget.SearchView
import androidx.recyclerview.widget.LinearLayoutManager
import com.google.android.material.snackbar.Snackbar
import id.my.mrgsrylm.githubuser.R
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import id.my.mrgsrylm.githubuser.databinding.ActivityMainBinding
import id.my.mrgsrylm.githubuser.view.adapter.ListUserAdapter
import id.my.mrgsrylm.githubuser.viewmodel.MainViewModel

class MainActivity : AppCompatActivity() {

    private val mainViewModel by viewModels<MainViewModel>()
    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        mainViewModel.isLoading.observe(this) {
            showLoading(it)
        }

        mainViewModel.user.observe(this) { user ->
            setUserData(user)
        }

        mainViewModel.snackBarText.observe(this) {
            Snackbar.make(
                window.decorView.rootView,
                it,
                Snackbar.LENGTH_SHORT
            ).show()
        }

        val layoutManager = LinearLayoutManager(this)
        binding.rvUser.layoutManager = layoutManager
    }

    override fun onCreateOptionsMenu(menu: Menu): Boolean {
        val inflater = menuInflater
        inflater.inflate(R.menu.option_menu, menu)

        val searchManager = getSystemService(Context.SEARCH_SERVICE) as SearchManager
        val searchView = menu.findItem(R.id.search).actionView as SearchView

        searchView.setSearchableInfo(searchManager.getSearchableInfo(componentName))
        searchView.queryHint = resources.getString(R.string.search_hint)
        searchView.setOnQueryTextListener(object : SearchView.OnQueryTextListener {
            override fun onQueryTextSubmit(query: String): Boolean {
                setNotFound(false)
                mainViewModel.findUser(query)
                searchView.clearFocus()
                return true
            }

            override fun onQueryTextChange(newText: String): Boolean {
                val emptyUser = listOf<UserItemResponse>()
                setNotFound(newText.isEmpty())
                setUserData(emptyUser)
                return false
            }
        })
        return true
    }

    private fun setUserData(githubUsername: List<UserItemResponse>) {
        val listUser = ArrayList<UserItemResponse>()
        for (i in githubUsername) {
            listUser.clear()
            listUser.addAll(githubUsername)
        }
        val adapter = ListUserAdapter(listUser)
        binding.rvUser.adapter = adapter

        adapter.setOnItemClickCallback(object
            : ListUserAdapter.OnItemClickCallback {
            override fun onItemClicked(data: UserItemResponse) {
                moveToDetail(data)
            }
        })
    }

    private fun moveToDetail(data: UserItemResponse) {
        val detailUserIntent = Intent(this@MainActivity, UserDetailActivity::class.java)
        detailUserIntent.putExtra(UserDetailActivity.EXTRA_USER, data.login)
        detailUserIntent.putExtra(UserDetailActivity.EXTRA_ACTIONBAR, data.login)
        startActivity(detailUserIntent)
    }

    private fun showLoading(isLoading: Boolean) {
        if (isLoading) {
            binding.progressBar.visibility = View.VISIBLE
        } else {
            binding.progressBar.visibility = View.GONE
        }
    }

    private fun setNotFound(notFound: Boolean) {
        if (notFound) {
            binding.tvNotFound.visibility = View.VISIBLE
        } else {
            binding.tvNotFound.visibility = View.GONE
        }
    }

    companion object {
        const val TAG = "MainActivity"
        const val GITHUB_USERNAME = "mrgsrylm"
    }
}