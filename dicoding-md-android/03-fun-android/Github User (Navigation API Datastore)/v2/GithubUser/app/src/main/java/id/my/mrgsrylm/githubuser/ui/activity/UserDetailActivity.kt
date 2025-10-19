package id.my.mrgsrylm.githubuser.ui.activity

import android.os.Bundle
import android.util.Log
import android.view.View
import androidx.activity.viewModels
import androidx.annotation.StringRes
import androidx.appcompat.app.AppCompatActivity
import com.bumptech.glide.Glide
import com.google.android.material.tabs.TabLayoutMediator
import id.my.mrgsrylm.githubuser.R
import id.my.mrgsrylm.githubuser.data.local.entity.FavoriteUser
import id.my.mrgsrylm.githubuser.data.remote.response.UserDetail
import id.my.mrgsrylm.githubuser.databinding.ActivityUserDetailBinding
import id.my.mrgsrylm.githubuser.helper.ViewModelFactory
import id.my.mrgsrylm.githubuser.ui.adapter.SectionsPagerAdapter
import id.my.mrgsrylm.githubuser.ui.viewmodel.UserDetailViewModel

class UserDetailActivity : AppCompatActivity() {

    private lateinit var binding: ActivityUserDetailBinding
    private val userDetailViewModel by viewModels<UserDetailViewModel> {
        ViewModelFactory.getInstance(application)
    }

    companion object {
        const val EXTRA_USER = "extra_user"
        const val EXTRA_FRAGMENT = "extra_fragment"

        @StringRes
        private val TAB_TITLES = intArrayOf(
            R.string.tab_followers,
            R.string.tab_following,
        )
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityUserDetailBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val username = intent.getStringExtra(EXTRA_USER)
        username?.let {
            userDetailViewModel.getUserDetail(it)
        }

        userDetailViewModel.user.observe(this) { user ->
            setUserDetailData(user)
            setFavoriteUser(user, userDetailViewModel)
        }

        userDetailViewModel.isLoading.observe(this) {
            showLoading(it)
        }

        setSectionsPager(userDetailViewModel)
    }

    private fun setUserDetailData(user: UserDetail) {
        binding.apply {
            Glide.with(this@UserDetailActivity)
                .load(user.avatarUrl)
                .into(ivDetailPhoto)
            tvDetailName.text = user.name
            tvDetailUsername.text = user.login
            tvFollowers.text = getString(R.string.followers, user.followers.toString())
            tvFollowing.text = getString(R.string.following, user.following.toString())
        }
    }

    private fun setSectionsPager(viewModel: UserDetailViewModel) {
        val userIntent = intent.extras
        if (userIntent != null) {
            val userLogin = userIntent.getString(EXTRA_USER)
            viewModel.getUserDetail(userLogin!!)

            val login = Bundle()
            login.putString(EXTRA_FRAGMENT, userLogin)

            val sectionsPagerAdapter = SectionsPagerAdapter(this, login)
            binding.viewPager.adapter = sectionsPagerAdapter
            TabLayoutMediator(binding.tabs, binding.viewPager) { tab, position ->
                tab.text = resources.getString(TAB_TITLES[position])
            }.attach()
            supportActionBar?.elevation = 0f
        }
    }

    private fun setFavoriteUser(user: UserDetail, viewModel: UserDetailViewModel) {
        viewModel.getFavoriteByUsername(user.login).observe(this) {
            val favoriteUser = FavoriteUser(user.login, user.avatarUrl)

            Log.d("UserDetailAct", favoriteUser.toString())
            var isFavorite = false

            if (it != null) {
                isFavorite = true
                binding.fabFavorite.setImageResource(R.drawable.baseline_love)
            } else {
                binding.fabFavorite.setImageResource(R.drawable.baseline_love_line)
            }

            binding.fabFavorite.setOnClickListener {
                if (!isFavorite) {
                    viewModel.insertFavorite(favoriteUser)
                    isFavorite = true
                    binding.fabFavorite.setImageResource(R.drawable.baseline_love)
                } else {
                    viewModel.deleteFavoriteByUsername(favoriteUser.username)
                    isFavorite = false
                    binding.fabFavorite.setImageResource(R.drawable.baseline_love_line)
                }
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
    }
}