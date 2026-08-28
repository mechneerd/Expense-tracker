package com.expensetracker.family

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.FragmentManager
import com.expensetracker.family.databinding.ActivityMainBinding
import com.google.android.material.bottomnavigation.BottomNavigationView

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        // Set up BottomNavigation listener
        binding.bottomNav.setOnItemSelectedListener { item ->
            when (item.itemId) {
                R.id.nav_home -> {
                    supportFragmentManager
                        .beginTransaction()
                        .replace(R.id.fragment_container, HomeFragment())
                        .commitNow()
                }
                R.id.nav_transactions -> {
                    supportFragmentManager
                        .beginTransaction()
                        .replace(R.id.fragment_container, TransactionEntryFragment())
                        .commitNow()
                }
                R.id.nav_categories -> {
                    supportFragmentManager
                        .beginTransaction()
                        .replace(R.id.fragment_container, CategoriesFragment())
                        .commitNow()
                }
                R.id.nav_dashboard -> {
                    // Dashboard only for family head
                    val userRole = android.preference.Preferences
                        .getDefaultSharedPreferences(this)
                        .getString("family_role", "")
                    if (userRole == "HEAD") {
                        supportFragmentManager
                            .beginTransaction()
                            .replace(R.id.fragment_container, DashboardFragment())
                            .commitNow()
                    } else {
                        android.widget.Toast.makeText(this,
                            "Dashboard only for Family Heads", android.widget.Toast.LENGTH_SHORT).show()
                        binding.bottomNav.selectedItemId = R.id.nav_home
                    }
                }
                R.id.nav_families -> {
                    supportFragmentManager
                        .beginTransaction()
                        .replace(R.id.fragment_container, FamilyCreationFragment())
                        .commitNow()
                }
                R.id.nav_profile -> {
                    supportFragmentManager
                        .beginTransaction()
                        .replace(R.id.fragment_container, ProfileFragment())
                        .commitNow()
                }
            }
            true
        }

        // Set default fragment (Home)
        supportFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, HomeFragment())
            .commitNow()
    }
}