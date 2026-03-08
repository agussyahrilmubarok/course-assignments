package id.my.mrgsrylm.githubuser.view.fragment

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import androidx.recyclerview.widget.LinearLayoutManager
import id.my.mrgsrylm.githubuser.data.remote.response.UserItemResponse
import id.my.mrgsrylm.githubuser.databinding.FragmentFollowersBinding
import id.my.mrgsrylm.githubuser.view.activity.UserDetailActivity
import id.my.mrgsrylm.githubuser.view.adapter.ListFollowersAdapter
import id.my.mrgsrylm.githubuser.viewmodel.FollowersViewModel

class FollowersFragment : Fragment() {

    private val followersViewModel by viewModels<FollowersViewModel>()
    private lateinit var binding: FragmentFollowersBinding

    companion object {
        const val TAG = "FragmentFollowers"
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentFollowersBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        followersViewModel.findFollowers(
            arguments?.getString(UserDetailActivity.EXTRA_FRAGMENT).toString()
        )

        followersViewModel.followers.observe(viewLifecycleOwner) { followers ->
            setFollowersList(followers)
        }

        followersViewModel.isLoading.observe(viewLifecycleOwner) {
            showLoading(it)
        }
    }

    private fun setFollowersList(followers: List<UserItemResponse>) {
        val users = ArrayList<UserItemResponse>()
        for (i in followers) {
            users.clear()
            users.addAll(followers)
        }
        binding.rvFollowers.layoutManager = LinearLayoutManager(requireActivity())
        val adapter = ListFollowersAdapter(users)
        binding.rvFollowers.adapter = adapter
    }

    private fun showLoading(isLoading: Boolean) {
        if (isLoading) {
            binding.progressBar.visibility = View.VISIBLE
        } else {
            binding.progressBar.visibility = View.GONE
        }
    }
}