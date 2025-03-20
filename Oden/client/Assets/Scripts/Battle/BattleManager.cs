using System;
using System.Collections;
using System.Collections.Generic;
using Oden.API;
using Oden.Hero;
using UnityEngine;

namespace Oden.Battle
{
    public class BattleManager : MonoBehaviour
    {
        // References to other managers
        private GameManager _gameManager;
        
        // Battle state
        private bool _isBattleInProgress = false;
        private BattleResult _currentBattleResult;
        private int _currentTurn = 0;
        private float _turnDelay = 1.0f; // Delay between turns for visualization
        
        // Events
        public event Action<BattleResult> OnBattleCompleted;
        public event Action<int> OnTurnCompleted;
        public event Action<BattleAction> OnActionPerformed;
        public event Action OnBattleStarted;
        
        private void Awake()
        {
            _gameManager = GameManager.Instance;
        }
        
        /// <summary>
        /// Start a battle with the current team
        /// </summary>
        public void StartBattle(string stageId)
        {
            if (_isBattleInProgress)
            {
                Debug.LogWarning("Battle already in progress");
                return;
            }
            
            _isBattleInProgress = true;
            _currentTurn = 0;
            
            // Notify battle started
            OnBattleStarted?.Invoke();
            
            // Call API to get battle result
            ApiClient.Instance.StartBattle(stageId, (success, result) =>
            {
                if (success && result != null)
                {
                    _currentBattleResult = result;
                    
                    // Start visualizing the battle
                    StartCoroutine(VisualizeBattle());
                }
                else
                {
                    Debug.LogError("Failed to start battle");
                    _isBattleInProgress = false;
                }
            });
        }
        
        /// <summary>
        /// Visualize the battle turn by turn
        /// </summary>
        private IEnumerator VisualizeBattle()
        {
            // Process each turn
            foreach (var turn in _currentBattleResult.battle_log)
            {
                _currentTurn = turn.turn;
                
                // Process each action in the turn
                foreach (var action in turn.actions)
                {
                    // Notify that an action was performed
                    OnActionPerformed?.Invoke(action);
                    
                    // Wait a short time to visualize the action
                    yield return new WaitForSeconds(_turnDelay / turn.actions.Count);
                }
                
                // Notify that a turn was completed
                OnTurnCompleted?.Invoke(_currentTurn);
                
                // Wait between turns
                yield return new WaitForSeconds(_turnDelay);
            }
            
            // Battle completed
            _isBattleInProgress = false;
            
            // Process rewards
            ProcessRewards();
            
            // Notify that the battle is completed
            OnBattleCompleted?.Invoke(_currentBattleResult);
        }
        
        /// <summary>
        /// Process battle rewards
        /// </summary>
        private void ProcessRewards()
        {
            if (_currentBattleResult == null || _currentBattleResult.rewards == null)
                return;
            
            // Add gold
            _gameManager.AddGold(_currentBattleResult.rewards.gold);
            
            // Add experience to heroes
            foreach (var expPair in _currentBattleResult.rewards.experience)
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
        /// Set the turn delay for battle visualization
        /// </summary>
        public void SetTurnDelay(float delay)
        {
            _turnDelay = Mathf.Clamp(delay, 0.1f, 5.0f);
        }
        
        /// <summary>
        /// Check if a battle is in progress
        /// </summary>
        public bool IsBattleInProgress()
        {
            return _isBattleInProgress;
        }
        
        /// <summary>
        /// Get the current battle result
        /// </summary>
        public BattleResult GetCurrentBattleResult()
        {
            return _currentBattleResult;
        }
        
        /// <summary>
        /// Get the current turn
        /// </summary>
        public int GetCurrentTurn()
        {
            return _currentTurn;
        }
    }
    
    /// <summary>
    /// Battle result data for client side processing
    /// </summary>
    public class BattleResult
    {
        public string battle_id;
        public string result;
        public List<BattleTurn> battle_log;
        public BattleRewards rewards;
        
        public bool IsVictory => result == "victory";
    }
    
    /// <summary>
    /// A turn in the battle
    /// </summary>
    public class BattleTurn
    {
        public int turn;
        public List<BattleAction> actions;
    }
    
    /// <summary>
    /// An action performed during a battle turn
    /// </summary>
    public class BattleAction
    {
        public string actor;        // ID of the acting entity (hero or enemy)
        public string target;       // ID of the target entity
        public string skill_used;   // ID of the skill used
        public int damage_dealt;    // Damage dealt by the action
        public int target_hp_remaining; // Remaining HP of the target after the action
        
        // Helper properties
        public bool IsHeroAction => actor.StartsWith("hero_");
        public bool IsEnemyAction => actor.StartsWith("enemy_");
        public bool IsTargetHero => target.StartsWith("hero_");
        public bool IsTargetEnemy => target.StartsWith("enemy_");
    }
    
    /// <summary>
    /// Rewards from a battle
    /// </summary>
    public class BattleRewards
    {
        public int gold;
        public Dictionary<string, int> experience;
        public List<string> items;
    }
} 