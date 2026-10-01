#include <algorithm>
#include <map>
#include <string>
#include <vector>
using namespace std;

vector<vector<string>> groupByLetters(vector<string>& words) {
    map<string, vector<string>> groups;
    for (const string& w : words) {
        string key = w;
        sort(key.begin(), key.end());
        groups[key].push_back(w);
    }
    vector<vector<string>> out;
    for (auto& [key, g] : groups) out.push_back(g);
    return out;
}
