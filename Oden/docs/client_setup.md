# Oden - Client Setup Guide

This guide will walk you through setting up the Unity client for Oden Idle RPG development.

## Prerequisites

- Unity 2021.3 LTS or newer
- Visual Studio 2019/2022 or Visual Studio Code with Unity extension
- Basic knowledge of C# and Unity
- Git client
- Xcode (for iOS builds)

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://your-repository-url/oden.git
cd oden
```

### 2. Install Unity

1. Download and install Unity Hub from [unity3d.com](https://unity3d.com/get-unity/download)
2. Install Unity 2021.3 LTS via Unity Hub
3. Add iOS Build Support module during installation

### 3. Open the Project

1. Launch Unity Hub
2. Click "Add" and browse to the `client` folder in your cloned repository
3. Select the project and open it with Unity 2021.3 LTS

### 4. Install Required Packages

Once the project is open, install the following packages via the Package Manager (Window > Package Manager):

- TextMeshPro
- Unity UI
- Newtonsoft Json
- Addressables (for asset management)
- Unity Test Framework (for testing)

### 5. Configure Project Settings

1. Go to Edit > Project Settings
2. Configure the following:
   - Player Settings:
     - Product Name: "Oden"
     - Company Name: Your company name
     - Version: 0.1.0
     - Bundle Identifier: com.yourcompany.oden
   - Quality Settings:
     - Set appropriate quality levels for mobile
   - Input System:
     - Configure touch input

## Project Structure

The client project is organized as follows:

```
client/
├── Assets/
│   ├── Animations/       # Animation clips and controllers
│   ├── Audio/            # Sound effects and music
│   ├── Prefabs/          # Reusable game objects
│   ├── Resources/        # Assets loaded at runtime
│   ├── Scenes/           # Unity scenes
│   │   ├── Login.unity   # Login and registration scene
│   │   ├── Main.unity    # Main hub scene
│   │   └── Battle.unity  # Battle scene
│   ├── Scripts/          # C# code
│   │   ├── API/          # Server communication
│   │   ├── Auth/         # Authentication
│   │   ├── Battle/       # Battle system
│   │   ├── Config/       # Configuration
│   │   ├── Hero/         # Hero management
│   │   ├── Idle/         # Idle/AFK system
│   │   ├── UI/           # User interface
│   │   └── Utils/        # Utility functions
│   ├── Textures/         # Images and sprites
│   └── UI/               # UI assets
└── ProjectSettings/      # Unity project settings
```

## Configuration

### Setting API Endpoint

Create a GameConfig script to manage your configuration:

1. Create a new C# script at `Assets/Scripts/Config/GameConfig.cs`:

```csharp
using UnityEngine;

namespace Oden.Config
{
    [CreateAssetMenu(fileName = "GameConfig", menuName = "Oden/Game Config")]
    public class GameConfig : ScriptableObject
    {
        [Header("API Settings")]
        [SerializeField] private string apiEndpoint = "http://localhost:8080/v1";
        
        [Header("Game Settings")]
        [SerializeField] private float idleRewardRate = 0.5f;  // resources per minute
        [SerializeField] private int maxTeamSize = 5;
        
        // Properties
        public string ApiEndpoint => apiEndpoint;
        public float IdleRewardRate => idleRewardRate;
        public int MaxTeamSize => maxTeamSize;
        
        // Singleton instance
        private static GameConfig _instance;
        public static GameConfig Instance
        {
            get
            {
                if (_instance == null)
                {
                    _instance = Resources.Load<GameConfig>("GameConfig");
                    if (_instance == null)
                    {
                        Debug.LogError("GameConfig not found in Resources folder!");
                    }
                }
                return _instance;
            }
        }
    }
}
```

2. Create the asset:
   - Right-click in the Project window
   - Select Create > Oden > Game Config
   - Name it "GameConfig"
   - Save it in the `Assets/Resources` folder

3. Configure the asset with appropriate values for your environment

## Creating Core Components

### 1. API Client

Create a script to handle API communication:

```csharp
// Assets/Scripts/API/ApiClient.cs
using System;
using System.Collections;
using System.Text;
using Newtonsoft.Json;
using Oden.Config;
using UnityEngine;
using UnityEngine.Networking;

namespace Oden.API
{
    public class ApiClient : MonoBehaviour
    {
        private static ApiClient _instance;
        public static ApiClient Instance
        {
            get
            {
                if (_instance == null)
                {
                    GameObject go = new GameObject("ApiClient");
                    _instance = go.AddComponent<ApiClient>();
                    DontDestroyOnLoad(go);
                }
                return _instance;
            }
        }

        private string _token;
        public bool IsAuthenticated => !string.IsNullOrEmpty(_token);

        public void SetToken(string token)
        {
            _token = token;
            PlayerPrefs.SetString("AuthToken", token);
            PlayerPrefs.Save();
        }

        private void Awake()
        {
            if (_instance != null && _instance != this)
            {
                Destroy(gameObject);
                return;
            }
            
            _instance = this;
            DontDestroyOnLoad(gameObject);
            
            // Try to load saved token
            _token = PlayerPrefs.GetString("AuthToken", "");
        }

        public void Register(string username, string email, string password, Action<bool, string> callback)
        {
            var registerData = new
            {
                username,
                email,
                password
            };
            
            StartCoroutine(PostRequest("auth/register", registerData, false, (success, response) =>
            {
                if (success)
                {
                    var responseObj = JsonConvert.DeserializeObject<RegisterResponse>(response);
                    SetToken(responseObj.token);
                    callback(true, "Registration successful");
                }
                else
                {
                    callback(false, response);
                }
            }));
        }

        public void Login(string username, string password, Action<bool, string> callback)
        {
            var loginData = new
            {
                username,
                password
            };
            
            StartCoroutine(PostRequest("auth/login", loginData, false, (success, response) =>
            {
                if (success)
                {
                    var responseObj = JsonConvert.DeserializeObject<LoginResponse>(response);
                    SetToken(responseObj.token);
                    callback(true, "Login successful");
                }
                else
                {
                    callback(false, response);
                }
            }));
        }

        public void GetHeroes(Action<bool, HeroList> callback)
        {
            StartCoroutine(GetRequest("heroes/list", true, (success, response) =>
            {
                if (success)
                {
                    var heroes = JsonConvert.DeserializeObject<HeroList>(response);
                    callback(true, heroes);
                }
                else
                {
                    callback(false, null);
                }
            }));
        }

        public void SummonHero(string summonType, Action<bool, SummonResult> callback)
        {
            var summonData = new
            {
                summon_type = summonType
            };
            
            StartCoroutine(PostRequest("heroes/summon", summonData, true, (success, response) =>
            {
                if (success)
                {
                    var result = JsonConvert.DeserializeObject<SummonResult>(response);
                    callback(true, result);
                }
                else
                {
                    callback(false, null);
                }
            }));
        }

        public void SaveTeam(Dictionary<int, string> positions, Action<bool, string> callback)
        {
            var teamData = new
            {
                positions
            };
            
            StartCoroutine(PostRequest("team/save", teamData, true, (success, response) =>
            {
                callback(success, response);
            }));
        }

        public void StartBattle(string stageId, Action<bool, BattleResult> callback)
        {
            var battleData = new
            {
                stage_id = stageId
            };
            
            StartCoroutine(PostRequest("battle/start", battleData, true, (success, response) =>
            {
                if (success)
                {
                    var result = JsonConvert.DeserializeObject<BattleResult>(response);
                    callback(true, result);
                }
                else
                {
                    callback(false, null);
                }
            }));
        }

        public void GetIdleRewards(Action<bool, IdleRewards> callback)
        {
            StartCoroutine(GetRequest("idle/rewards", true, (success, response) =>
            {
                if (success)
                {
                    var rewards = JsonConvert.DeserializeObject<IdleRewards>(response);
                    callback(true, rewards);
                }
                else
                {
                    callback(false, null);
                }
            }));
        }

        public void ClaimIdleRewards(Action<bool, IdleRewards> callback)
        {
            StartCoroutine(PostRequest("idle/claim", null, true, (success, response) =>
            {
                if (success)
                {
                    var rewards = JsonConvert.DeserializeObject<IdleRewards>(response);
                    callback(true, rewards);
                }
                else
                {
                    callback(false, null);
                }
            }));
        }

        private IEnumerator GetRequest(string endpoint, bool requireAuth, Action<bool, string> callback)
        {
            string url = $"{GameConfig.Instance.ApiEndpoint}/{endpoint}";
            
            using (UnityWebRequest www = UnityWebRequest.Get(url))
            {
                if (requireAuth && !string.IsNullOrEmpty(_token))
                {
                    www.SetRequestHeader("Authorization", $"Bearer {_token}");
                }
                
                yield return www.SendWebRequest();
                
                if (www.result == UnityWebRequest.Result.Success)
                {
                    callback(true, www.downloadHandler.text);
                }
                else
                {
                    Debug.LogError($"API Error: {www.error}");
                    callback(false, www.error);
                }
            }
        }

        private IEnumerator PostRequest(string endpoint, object data, bool requireAuth, Action<bool, string> callback)
        {
            string url = $"{GameConfig.Instance.ApiEndpoint}/{endpoint}";
            string jsonData = data != null ? JsonConvert.SerializeObject(data) : "{}";
            
            using (UnityWebRequest www = new UnityWebRequest(url, "POST"))
            {
                byte[] bodyRaw = Encoding.UTF8.GetBytes(jsonData);
                www.uploadHandler = new UploadHandlerRaw(bodyRaw);
                www.downloadHandler = new DownloadHandlerBuffer();
                www.SetRequestHeader("Content-Type", "application/json");
                
                if (requireAuth && !string.IsNullOrEmpty(_token))
                {
                    www.SetRequestHeader("Authorization", $"Bearer {_token}");
                }
                
                yield return www.SendWebRequest();
                
                if (www.result == UnityWebRequest.Result.Success)
                {
                    callback(true, www.downloadHandler.text);
                }
                else
                {
                    Debug.LogError($"API Error: {www.error}");
                    callback(false, www.error);
                }
            }
        }
    }

    // Response models
    [Serializable]
    public class LoginResponse
    {
        public bool success;
        public string user_id;
        public string token;
    }

    [Serializable]
    public class RegisterResponse
    {
        public bool success;
        public string user_id;
        public string token;
    }

    [Serializable]
    public class HeroList
    {
        public List<Hero> heroes;
    }

    [Serializable]
    public class Hero
    {
        public string id;
        public string hero_type_id;
        public string name;
        public int level;
        public int experience;
        public int hp;
        public int atk;
        public List<Skill> skills;
    }

    [Serializable]
    public class Skill
    {
        public string id;
        public string name;
        public string description;
        public float damage_multiplier;
        public int cooldown;
    }

    [Serializable]
    public class SummonResult
    {
        public bool success;
        public List<Hero> heroes;
    }

    [Serializable]
    public class BattleResult
    {
        public string battle_id;
        public string result;
        public List<BattleTurn> battle_log;
        public BattleRewards rewards;
    }

    [Serializable]
    public class BattleTurn
    {
        public int turn;
        public List<BattleAction> actions;
    }

    [Serializable]
    public class BattleAction
    {
        public string actor;
        public string target;
        public string skill_used;
        public int damage_dealt;
        public int target_hp_remaining;
    }

    [Serializable]
    public class BattleRewards
    {
        public int gold;
        public Dictionary<string, int> experience;
        public List<string> items;
    }

    [Serializable]
    public class IdleRewards
    {
        public int time_away;
        public IdleRewardData rewards;
    }

    [Serializable]
    public class IdleRewardData
    {
        public int gold;
        public Dictionary<string, int> experience;
    }
}
```

### 2. Game Manager

Create a central manager for game state:

```csharp
// Assets/Scripts/GameManager.cs
using System;
using Oden.API;
using Oden.Hero;
using UnityEngine;
using UnityEngine.SceneManagement;

namespace Oden
{
    public class GameManager : MonoBehaviour
    {
        private static GameManager _instance;
        public static GameManager Instance
        {
            get
            {
                if (_instance == null)
                {
                    GameObject go = new GameObject("GameManager");
                    _instance = go.AddComponent<GameManager>();
                    DontDestroyOnLoad(go);
                }
                return _instance;
            }
        }

        // Player data
        public int Gold { get; private set; }
        public HeroManager HeroManager { get; private set; }
        
        // Events
        public event Action<int> OnGoldChanged;
        public event Action OnPlayerDataLoaded;

        private void Awake()
        {
            if (_instance != null && _instance != this)
            {
                Destroy(gameObject);
                return;
            }
            
            _instance = this;
            DontDestroyOnLoad(gameObject);
            
            // Initialize components
            HeroManager = gameObject.AddComponent<HeroManager>();
        }

        private void Start()
        {
            // Check if we're authenticated
            if (ApiClient.Instance.IsAuthenticated)
            {
                LoadPlayerData();
            }
            else
            {
                // Go to login scene if not authenticated
                SceneManager.LoadScene("Login");
            }
        }

        public void LoadPlayerData()
        {
            // Load player's heroes
            HeroManager.LoadHeroes(() =>
            {
                // More data loading can happen here
                Gold = 1000; // Placeholder, should come from API
                OnGoldChanged?.Invoke(Gold);
                OnPlayerDataLoaded?.Invoke();
            });
        }

        public void AddGold(int amount)
        {
            Gold += amount;
            OnGoldChanged?.Invoke(Gold);
        }

        public bool SpendGold(int amount)
        {
            if (Gold >= amount)
            {
                Gold -= amount;
                OnGoldChanged?.Invoke(Gold);
                return true;
            }
            return false;
        }
    }
}
```

## Next Steps

After setting up the project structure and core components:

1. Create UI screens for:
   - Login/Registration
   - Main hub
   - Hero collection
   - Team formation
   - Battle screen
   - Summon heroes screen

2. Implement the HeroManager and other game systems

3. Follow the [Testing Guide](testing.md) to test your implementation

4. Refer to the [Deployment Guide](deployment.md) when ready to build for iOS

## Troubleshooting

- **Unity Package Errors**: Try clearing the Library folder and reopening the project
- **iOS Build Issues**: Make sure Xcode and iOS modules are up to date
- **API Connection Issues**: Verify server is running and API endpoint is correctly configured 