using System;
using System.Collections.Generic;
using System.Linq;
using Oden.API;
using UnityEngine;

namespace Oden.Hero
{
    public class HeroManager : MonoBehaviour
    {
        // Events
        public event Action OnHeroesLoaded;
        public event Action<Hero> OnHeroAdded;
        public event Action<string> OnTeamSaved;
        
        // Collection of player's heroes
        private Dictionary<string, Hero> _heroes = new Dictionary<string, Hero>();
        public IReadOnlyDictionary<string, Hero> Heroes => _heroes;
        
        // Current team formation
        private Dictionary<int, string> _teamFormation = new Dictionary<int, string>();
        public IReadOnlyDictionary<int, string> TeamFormation => _teamFormation;
        
        // Max team size
        private const int MaxTeamSize = 5;
        
        /// <summary>
        /// Load heroes from the server
        /// </summary>
        public void LoadHeroes(Action onComplete = null)
        {
            ApiClient.Instance.GetHeroes((success, heroList) =>
            {
                if (success && heroList != null)
                {
                    _heroes.Clear();
                    foreach (var hero in heroList.heroes)
                    {
                        _heroes[hero.id] = new Hero(hero);
                    }
                    
                    Debug.Log($"Loaded {_heroes.Count} heroes");
                    OnHeroesLoaded?.Invoke();
                }
                else
                {
                    Debug.LogError("Failed to load heroes");
                }
                
                onComplete?.Invoke();
            });
        }
        
        /// <summary>
        /// Load team formation from the server
        /// </summary>
        public void LoadTeamFormation(Action<bool> callback = null)
        {
            ApiClient.Instance.GetTeam((success, teamData) =>
            {
                if (success && teamData != null)
                {
                    _teamFormation.Clear();
                    foreach (var position in teamData.positions)
                    {
                        int positionIndex = int.Parse(position.Key);
                        if (position.Value != null)
                        {
                            string heroId = position.Value.hero_id;
                            _teamFormation[positionIndex] = heroId;
                        }
                    }
                    
                    Debug.Log("Team formation loaded");
                    callback?.Invoke(true);
                }
                else
                {
                    Debug.LogError("Failed to load team formation");
                    callback?.Invoke(false);
                }
            });
        }
        
        /// <summary>
        /// Save team formation to the server
        /// </summary>
        public void SaveTeamFormation(Action<bool> callback = null)
        {
            ApiClient.Instance.SaveTeam(_teamFormation, (success, response) =>
            {
                if (success)
                {
                    Debug.Log("Team formation saved");
                    OnTeamSaved?.Invoke(response);
                    callback?.Invoke(true);
                }
                else
                {
                    Debug.LogError($"Failed to save team formation: {response}");
                    callback?.Invoke(false);
                }
            });
        }
        
        /// <summary>
        /// Set a hero in a team position
        /// </summary>
        public bool SetTeamPosition(int position, string heroId)
        {
            // Validate position
            if (position < 1 || position > MaxTeamSize)
            {
                Debug.LogError($"Invalid team position: {position}");
                return false;
            }
            
            // Validate hero
            if (!string.IsNullOrEmpty(heroId) && !_heroes.ContainsKey(heroId))
            {
                Debug.LogError($"Hero not found: {heroId}");
                return false;
            }
            
            // If hero is already in another position, remove it
            if (!string.IsNullOrEmpty(heroId))
            {
                foreach (var pos in _teamFormation.Keys.ToList())
                {
                    if (_teamFormation[pos] == heroId)
                    {
                        _teamFormation.Remove(pos);
                    }
                }
            }
            
            // Set or clear the position
            if (string.IsNullOrEmpty(heroId))
            {
                _teamFormation.Remove(position);
            }
            else
            {
                _teamFormation[position] = heroId;
            }
            
            return true;
        }
        
        /// <summary>
        /// Summon a new hero
        /// </summary>
        public void SummonHero(string summonType, Action<bool, List<Hero>> callback)
        {
            ApiClient.Instance.SummonHero(summonType, (success, result) =>
            {
                if (success && result != null && result.success)
                {
                    List<Hero> newHeroes = new List<Hero>();
                    foreach (var heroData in result.heroes)
                    {
                        Hero hero = new Hero(heroData);
                        _heroes[hero.Id] = hero;
                        newHeroes.Add(hero);
                        OnHeroAdded?.Invoke(hero);
                    }
                    
                    Debug.Log($"Summoned {newHeroes.Count} new heroes");
                    callback?.Invoke(true, newHeroes);
                }
                else
                {
                    Debug.LogError("Failed to summon hero");
                    callback?.Invoke(false, null);
                }
            });
        }
        
        /// <summary>
        /// Get heroes by filter
        /// </summary>
        public List<Hero> GetHeroesByFilter(Func<Hero, bool> filter)
        {
            return _heroes.Values.Where(filter).ToList();
        }
        
        /// <summary>
        /// Get highest level heroes
        /// </summary>
        public List<Hero> GetHighestLevelHeroes(int count)
        {
            return _heroes.Values.OrderByDescending(h => h.Level).Take(count).ToList();
        }
        
        /// <summary>
        /// Get a hero by ID
        /// </summary>
        public Hero GetHero(string id)
        {
            return _heroes.TryGetValue(id, out var hero) ? hero : null;
        }
        
        /// <summary>
        /// Get heroes in the current team
        /// </summary>
        public List<Hero> GetTeamHeroes()
        {
            List<Hero> team = new List<Hero>();
            
            foreach (var position in _teamFormation.Keys.OrderBy(p => p))
            {
                string heroId = _teamFormation[position];
                if (_heroes.TryGetValue(heroId, out var hero))
                {
                    team.Add(hero);
                }
            }
            
            return team;
        }
    }

    /// <summary>
    /// Represents a hero in the game
    /// </summary>
    public class Hero
    {
        // Basic properties
        public string Id { get; }
        public string HeroTypeId { get; }
        public string Name { get; }
        public int Level { get; private set; }
        public int Experience { get; private set; }
        public int HP { get; private set; }
        public int ATK { get; private set; }
        
        // Skills
        public List<Skill> Skills { get; }
        
        public Hero(ApiClient.Hero heroData)
        {
            Id = heroData.id;
            HeroTypeId = heroData.hero_type_id;
            Name = heroData.name;
            Level = heroData.level;
            Experience = heroData.experience;
            HP = heroData.hp;
            ATK = heroData.atk;
            
            Skills = new List<Skill>();
            if (heroData.skills != null)
            {
                foreach (var skillData in heroData.skills)
                {
                    Skills.Add(new Skill(skillData));
                }
            }
        }
        
        /// <summary>
        /// Add experience to the hero
        /// </summary>
        public void AddExperience(int amount)
        {
            Experience += amount;
            
            // Simple level up calculation (can be refined)
            int newLevel = 1 + (Experience / 100);
            if (newLevel > Level)
            {
                Level = newLevel;
                
                // Recalculate stats based on level
                HP = CalculateHP();
                ATK = CalculateATK();
            }
        }
        
        /// <summary>
        /// Calculate HP based on level and base stats
        /// </summary>
        private int CalculateHP()
        {
            // Simple calculation for MVP
            return HP * (1 + (Level - 1) * 0.1);
        }
        
        /// <summary>
        /// Calculate ATK based on level and base stats
        /// </summary>
        private int CalculateATK()
        {
            // Simple calculation for MVP
            return ATK * (1 + (Level - 1) * 0.1);
        }
    }

    /// <summary>
    /// Represents a hero skill
    /// </summary>
    public class Skill
    {
        public string Id { get; }
        public string Name { get; }
        public string Description { get; }
        public float DamageMultiplier { get; }
        public int Cooldown { get; }
        
        public Skill(ApiClient.Skill skillData)
        {
            Id = skillData.id;
            Name = skillData.name;
            Description = skillData.description;
            DamageMultiplier = skillData.damage_multiplier;
            Cooldown = skillData.cooldown;
        }
        
        /// <summary>
        /// Calculate damage dealt by this skill
        /// </summary>
        public int CalculateDamage(int baseAttack)
        {
            return Mathf.RoundToInt(baseAttack * DamageMultiplier);
        }
    }
} 