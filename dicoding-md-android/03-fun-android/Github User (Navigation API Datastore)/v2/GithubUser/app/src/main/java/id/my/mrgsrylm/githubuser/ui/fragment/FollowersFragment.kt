package id.my.mrgsrylm.githubuser.ui.fragment

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import androidx.recyclerview.widget.LinearLayoutManager
import id.my.mrgsrylm.githubuser.data.remote.response.UserItem
import id.my.mrgsrylm.githubuser.databinding.FragmentFollowersBinding
import id.my.mrgsrylm.githubuser.ui.activity.UserDetailActivity
import id.my.mrgsrylm.githubuser.ui.adapter.ItemUserAdapter
import id.my.mrgsrylm.githubuser.ui.viewmodel.FollowersViewModel

class FollowersFragment : Fragment() {

    private lateinit var binding: FragmentFollowersBinding
    private lateinit var adapter: ItemUserAdapter

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // Inflate the layout for this fragment
        binding = FragmentFollowersBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        adapter = ItemUserAdapter()
        binding.rvFollowers.adapter = adapter

        val layoutManager = LinearLayoutManager(requireContext())
        binding.rvFollowers.layoutManager = layoutManager

        val followersViewModel = ViewModelProvider(
            this,
            ViewModelProvider.NewInstanceFactory()
        ).get(FollowersViewModel::class.java)

        followersViewModel.getFollowers(
            arguments?.getString(UserDetailActivity.EXTRA_FRAGMENT).toString()
        )

        followersViewModel.followers.observe(viewLifecycleOwner) { followers ->
            setFollowersList(followers)
        }

        followersViewModel.isLoading.observe(viewLifecycleOwner) {
            showLoading(it)
        }
    }

    private fun setFollowersList(followers: List<UserItem>) {
        adapter.submitList(followers)
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
    }
}