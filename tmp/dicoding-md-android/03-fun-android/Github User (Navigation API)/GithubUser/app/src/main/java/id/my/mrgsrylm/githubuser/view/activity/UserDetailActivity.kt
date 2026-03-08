package id.my.mrgsrylm.githubuser.view.activity

import android.os.Bundle
import android.view.MenuItem
import android.view.View
import androidx.activity.viewModels
import androidx.annotation.StringRes
import androidx.appcompat.app.ActionBar
import androidx.appcompat.app.AppCompatActivity
import com.bumptech.glide.Glide
import com.google.android.material.snackbar.Snackbar
import com.google.android.material.tabs.TabLayoutMediator
import id.my.mrgsrylm.githubuser.R
import id.my.mrgsrylm.githubuser.data.remote.response.UserDetailResponse
import id.my.mrgsrylm.githubuser.databinding.ActivityUserDetailBinding
import id.my.mrgsrylm.githubuser.view.adapter.UserPagerAdapter
import id.my.mrgsrylm.githubuser.viewmodel.UserDetailViewModel

class UserDetailActivity : AppCompatActivity() {

    private val userDetailViewModel by viewModels<UserDetailViewModel>()
    private lateinit var binding: ActivityUserDetailBinding
    private lateinit var actionBar: ActionBar

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityUserDetailBinding.inflate(layoutInflater)
        setContentView(binding.root)

        actionBar = supportActionBar!!
        actionBar.title = intent.getStringExtra(EXTRA_ACTIONBAR).toString()
        actionBar.setDisplayHomeAsUpEnabled(true)

        userDetailViewModel.detail.observe(this) { userDetail ->
            setUserData(userDetail)
        }

        userDetailViewModel.isLoading.observe(this) {
            showLoading(it)
        }

        userDetailViewModel.snackBarText.observe(this) {
            Snackbar.make(
                window.decorView.rootView,
                it,
                Snackbar.LENGTH_SHORT
            ).show()
        }

        setFragment()

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

    private fun setFragment() {
        val userIntent = intent.extras
        if (userIntent != null) {
            val userLogin = userIntent.getString(EXTRA_USER)
            userDetailViewModel.getDetail(userLogin!!)

            val login = Bundle()
            login.putString(EXTRA_FRAGMENT, userLogin)

            val sectionsPagerAdapter = UserPagerAdapter(this, login)

            binding.viewPager.adapter = sectionsPagerAdapter
            TabLayoutMediator(binding.tabs, binding.viewPager) { tab, position ->
                tab.text = resources.getString(TAB_TITLES[position])
            }.attach()
            supportActionBar?.elevation = 0f
        }
    }

    private fun setUserData(username: UserDetailResponse) {
        binding.apply {
            Glide.with(this@UserDetailActivity)
                .load(username.avatarUrl)
                .into(ivDetailPhoto)
            tvDetailUsername.text = username.login
            if (username.name.toString() != "null")
                tvDetailName.text = username.name.toString()
            tvFollowers.text = getString(R.string.followers, username.followers.toString())
            tvFollowing.text = getString(R.string.following, username.following.toString())
        }
    }

    private fun showLoading(isLoading: Boolean) {
        if (isLoading) {
            binding.progressBar.visibility = View.VISIBLE
        } else {
            binding.progressBar.visibility = View.GONE
        }
    }

    companion object {
        const val EXTRA_ACTIONBAR = "extra_actionBar"
        const val EXTRA_USER = "extra_user"
        const val EXTRA_FRAGMENT = "extra_fragment"

        @StringRes
        val TAB_TITLES = intArrayOf(
            R.string.tab_followers,
            R.string.tab_following
        )
    }
}