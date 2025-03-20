using System;
using System.Collections;
using Oden.API;
using UnityEngine;

namespace Oden.Idle
{
    public class IdleManager : MonoBehaviour
    {
        // References to other managers
        private GameManager _gameManager;
        
        // Idle rewards data
        private IdleRewards _pendingRewards;
        private DateTime _lastCheckTime;
        
        // Check interval in seconds
        private float _checkInterval = 60.0f;
        
        // Events
        public event Action<IdleRewards> OnIdleRewardsReceived;
        public event Action<IdleRewards> OnIdleRewardsClaimed;
        
        private void Awake()
        {
            _gameManager = GameManager.Instance;
            _lastCheckTime = DateTime.Now;
        }
        
        private void Start()
        {
            // Check for idle rewards on startup
            CheckIdleRewards();
            
            // Start periodic checking
            StartCoroutine(PeriodicCheck());
        }
        
        /// <summary>
        /// Check for idle rewards periodically
        /// </summary>
        private IEnumerator PeriodicCheck()
        {
            while (true)
            {
                yield return new WaitForSeconds(_checkInterval);
                CheckIdleRewards();
            }
        }
        
        /// <summary>
        /// Check for idle rewards
        /// </summary>
        public void CheckIdleRewards()
        {
            ApiClient.Instance.GetIdleRewards((success, rewards) =>
            {
                if (success && rewards != null)
                {
                    _pendingRewards = rewards;
                    _lastCheckTime = DateTime.Now;
                    
                    // Notify that idle rewards are received
                    if (_pendingRewards.time_away > 0)
                    {
                        OnIdleRewardsReceived?.Invoke(_pendingRewards);
                    }
                }
                else
                {
                    Debug.LogError("Failed to check idle rewards");
                }
            });
        }
        
        /// <summary>
        /// Claim idle rewards
        /// </summary>
        public void ClaimIdleRewards()
        {
            if (_pendingRewards == null || _pendingRewards.time_away <= 0)
            {
                Debug.LogWarning("No idle rewards to claim");
                return;
            }
            
            ApiClient.Instance.ClaimIdleRewards((success, rewards) =>
            {
                if (success && rewards != null)
                {
                    // Process the rewards
                    ProcessRewards(rewards);
                    
                    // Clear pending rewards
                    _pendingRewards = null;
                    
                    // Notify that idle rewards are claimed
                    OnIdleRewardsClaimed?.Invoke(rewards);
                }
                else
                {
                    Debug.LogError("Failed to claim idle rewards");
                }
            });
        }
        
        /// <summary>
        /// Process idle rewards
        /// </summary>
        private void ProcessRewards(IdleRewards rewards)
        {
            if (rewards == null || rewards.rewards == null)
                return;
            
            // Add gold
            _gameManager.AddGold(rewards.rewards.gold);
            
            // Add experience to heroes
            foreach (var expPair in rewards.rewards.experience)
            {
                string heroId = expPair.Key;
                int expAmount = expPair.Value;
                
                Hero.Hero hero = _gameManager.HeroManager.GetHero(heroId);
                if (hero != null)
                {
                    hero.AddExperience(expAmount);
                }
            }
        }
        
        /// <summary>
        /// Get elapsed time since last check
        /// </summary>
        public TimeSpan GetElapsedTimeSinceLastCheck()
        {
            return DateTime.Now - _lastCheckTime;
        }
        
        /// <summary>
        /// Get pending rewards
        /// </summary>
        public IdleRewards GetPendingRewards()
        {
            return _pendingRewards;
        }
        
        /// <summary>
        /// Check if there are pending rewards
        /// </summary>
        public bool HasPendingRewards()
        {
            return _pendingRewards != null && _pendingRewards.time_away > 0;
        }
        
        /// <summary>
        /// Set check interval
        /// </summary>
        public void SetCheckInterval(float interval)
        {
            _checkInterval = Mathf.Max(10.0f, interval);
        }
    }
    
    /// <summary>
    /// Idle rewards data for client side processing
    /// </summary>
    public class IdleRewards
    {
        public int time_away;  // Time away in seconds
        public IdleRewardData rewards;
        
        // Helper properties
        public TimeSpan TimeAway => TimeSpan.FromSeconds(time_away);
        public string FormattedTimeAway
        {
            get
            {
                TimeSpan time = TimeAway;
                if (time.TotalHours >= 1)
                {
                    return $"{(int)time.TotalHours}h {time.Minutes}m";
                }
                else if (time.TotalMinutes >= 1)
                {
                    return $"{time.Minutes}m {time.Seconds}s";
                }
                else
                {
                    return $"{time.Seconds}s";
                }
            }
        }
    }
    
    /// <summary>
    /// Idle reward data
    /// </summary>
    public class IdleRewardData
    {
        public int gold;
        public System.Collections.Generic.Dictionary<string, int> experience;
    }
} 