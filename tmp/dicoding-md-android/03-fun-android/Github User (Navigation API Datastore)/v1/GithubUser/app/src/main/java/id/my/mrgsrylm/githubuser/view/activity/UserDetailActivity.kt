package id.my.mrgsrylm.githubuser.view.activity

import android.content.Intent
import android.os.Bundle
import android.view.Menu
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
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser
import id.my.mrgsrylm.githubuser.data.remote.response.UserDetailResponse
import id.my.mrgsrylm.githubuser.databinding.ActivityUserDetailBinding
import id.my.mrgsrylm.githubuser.view.adapter.UserPagerAdapter
import id.my.mrgsrylm.githubuser.viewmodel.ViewModelFactory
import id.my.mrgsrylm.githubuser.viewmodel.user.UserDetailViewModel

class UserDetailActivity : AppCompatActivity() {

    private lateinit var binding: ActivityUserDetailBinding
    private lateinit var actionBar: ActionBar
    private val userDetailViewModel by viewModels<UserDetailViewModel> {
        ViewModelFactory.getInstance(application)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityUserDetailBinding.inflate(layoutInflater)
        setContentView(binding.root)

        actionBar = supportActionBar!!
        actionBar.title = intent.getStringExtra(EXTRA_ACTIONBAR).toString()
        actionBar.setDisplayHomeAsUpEnabled(true)

        userDetailViewModel.detail.observe(this) { userDetail ->
            setUserData(userDetail)
            setUpFavoriteUser(userDetail)
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

    private fun setUpFavoriteUser(userDetail: UserDetailResponse) {
        userDetailViewModel.getFavoriteUser(userDetail.login ?: "").observe(this) {
            val favoriteUser = FavoriteUser(userDetail.login ?: "", userDetail.avatarUrl)
            var isFavorite = false

            if (it != null) {
                isFavorite = true
                binding.fabFavorite.setImageResource(R.drawable.baseline_love)
            }

            binding.fabFavorite.setOnClickListener {
                if (!isFavorite) {
                    userDetailViewModel.insert(favoriteUser)
                    isFavorite = true
                    binding.fabFavorite.setImageResource(R.drawable.baseline_love)
                } else {
                    userDetailViewModel.delete(favoriteUser.username)
                    isFavorite = false
                    binding.fabFavorite.setImageResource(R.drawable.baseline_love_line)
                }
            }
        }
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

    override fun onCreateOptionsMenu(menu: Menu?): Boolean {
        menuInflater.inflate(R.menu.share_menu, menu)
        return super.onCreateOptionsMenu(menu)
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        when (item.itemId) {
            android.R.id.home -> {
                finish()
                return true
            }

            R.id.action_share -> {
                userDetailViewModel.detail.observe(this) { userDetail ->
                    val shareIntent = Intent()
                    shareIntent.action = Intent.ACTION_SEND
                    shareIntent.putExtra(
                        Intent.EXTRA_TEXT,
                        "Hey, check out this user! \n${userDetail.htmlUrl}"
                    )
                    shareIntent.type = "text/plain"
                    startActivity(Intent.createChooser(shareIntent, "Share to: "))
                }
            }
        }
        return super.onContextItemSelected(item)
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