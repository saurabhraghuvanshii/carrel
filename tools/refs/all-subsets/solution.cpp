#include <string>
#include <vector>
using namespace std;

vector<string> subsets(string& word) {
    vector<string> out;
    int n = word.size();
    for (int mask = 0; mask < (1 << n); mask++) {
        string s;
        for (int i = 0; i < n; i++)
            if (mask >> i & 1) s += word[i];
        out.push_back(s);
    }
    return out;
}
