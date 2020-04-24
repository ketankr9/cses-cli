# cses-cli ![Go](https://github.com/ketankr9/cses-cli/workflows/Go/badge.svg)
solve cses.fi problemset via command line. Inspired from awesome [leetcode-cli](https://github.com/skygragon/leetcode-cli).


## How to use cses-cli
* Download the required binary from [https://github.com/ketankr9/cses-cli/releases](https://github.com/ketankr9/cses-cli/releases)
* Install [lynx](https://www.google.com/search?q=install+lynx+&oq=install+lynx). For ubuntu use ```sudo apt install lynx```
* Rename the binary to cses-cli and move it to PATH
* Just type these command in the terminal one by one and you will understand how to use it.
```
cses-cli login
cses-cli list
cses-cli show 1742
cses-cli solve 1742
cses-cli submit 1742.Robot-Path.cpp
```
## Auto Commit to your GitHub Repository (if needed)  
* Obtain access token for your repo. See [how to get token](https://github.com/skygragon/leetcode-cli-plugins/blob/master/docs/github.md#generate-token)    
* **Note:** Please make sure that your repo has *at least one* commit.  
* Configure cses-cli for github as follows.  
```
$$$ cses-cli github
Token: 1f10d6065e78a2654a14xxxxxxxxxxxxxxxxxxxx 
Repository: cses-solutions 
Github Username: ketankr9 
Github Email: ketan.krishnan.xxxxx@iitbhu.ac.in
```

## Configuration  
**Custom Template**: Create a file named *template.cpp, template.py, template.java, template.js* for respective languages in the **working** directory to use this feature.  
**Languages Supported**: *cpp*(C++17) \[default\], *java*, *python*(CPython3), *javascript*(Node.js)  
Also you can manually edit the config at ```~/.cses/config.json```  
```
$$$ cat template.cpp 
#include<bits/stdc++.h>
using namespace std;

#define pr pair<int, int>
#define f first
#define s second
#define mk make_pair

int main(){

	return 0;
}

$$$ cat ~/.cses/config.json 
{
    "csrf": "3e236b578af1e020070281xxxxxxxxxx",
    "username": "test123xyz",
    "cookie": "PHPSESSID=e17c5cae7fcxxxxxxxxxa7d14ee120fc4fb859c; path=/",
    "root": "/home/username/.cses",
    "editor": "atom",
    "lang": "java",
    "github": {
        "token": "1f10d6065e78a2654a14xxxxxxxxxxxxxxxxxxxx",
        "repository": "cses-solutions",
        "username": "ketankr9",
        "email": "ketan.krishnan.xxxxx@iitbhu.ac.in"
    }
}
```

## Screenshot
```
$$$ cses-cli login
Username: test123xyz
Password: [<DELETED>]
Logged in successfully

$$$ cses-cli list
	✔ [1068] Weird Algorithm           (95.6 %)
	✘ [1083] Missing Number            (92.1 %)
	- [1069] Repetitions               (93.9 %)
	- [1094] Increasing Array          (96.0 %)
	- [1070] Permutations              (96.4 %)
	- [1071] Number Spiral             (93.1 %)
	- [1072] Two Knights               (93.3 %)
	- [1092] Two Sets                  (94.1 %)
	- [1617] Bit Strings               (96.0 %)
	- [1618] Trailing Zeros            (94.1 %)
  [<DELETED>]
  
$$$ cses-cli show 1068
   CSES - Weird Algorithm
     * Time limit: 1.00 s
     * Memory limit: 512 MB

   Consider an algorithm that takes as input a positive integer $n$. If
   $n$ is even, the algorithm divides it by two, and if $n$ is odd, the
   algorithm multiplies it by three and adds one. The algorithm repeats
   this, until $n$ is one. For example, the sequence for $n=3$ is as
   follows:
   [ 3 → 10 → 5 → 16 → 8
   → 4 → 2 → 1]
   Your task is to simulate the execution of the algorithm for a given
   value of $n$.
   Input
   The only input line contains an integer $n$.
   Output
   Print a line that contains all values of $n$ during the algorithm.
   Constraints
     * $1 ≤ n ≤ 10^6$

   Example
   Input:
   3
   Output:
   3 10 5 16 8 4 2 1
   
//below command also opens editor with problem statement and code stub
$$$ cses-cli solve 1068
   CSES - Weird Algorithm
     * Time limit: 1.00 s
     * Memory limit: 512 MB

   Consider an algorithm that takes as input a positive integer $n$. If
   $n$ is even, the algorithm divides it by two, and if $n$ is odd, the
   algorithm multiplies it by three and adds one. The algorithm repeats
   this, until $n$ is one. For example, the sequence for $n=3$ is as
   follows:
   [ 3 → 10 → 5 → 16 → 8
   → 4 → 2 → 1]
   Your task is to simulate the execution of the algorithm for a given
   value of $n$.
   Input
   The only input line contains an integer $n$.
   Output
   Print a line that contains all values of $n$ during the algorithm.
   Constraints
     * $1 ≤ n ≤ 10^6$

   Example
   Input:
   3
   Output:
   3 10 5 16 8 4 2 1

$$$ cses-cli submit 1068.Weird-Algorithm.cpp 
Task:Weird Algorithm
Sender:test123xyz
Submission time:2020-03-07 13:56:29
Language:C++17
Status:READY
Result:ACCEPTED
Github: cses-solutions ✔
```

>I will add support for these features only if people show some love to this repo since current commit suffices my need.

*	Supports only C++ currently, will add support for other languages on request. **UPDATE: Done**
*	A modifiable template code file. **UPDATE: Done**
*	Auto commit to Github repository **UPDATE: Done**
*	Clean LaTex based $equation$ from problem statement. Eg $1 ≤ n ≤ 10^6$
