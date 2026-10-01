#include <algorithm>
#include <string>
#include <vector>
using namespace std;

vector<string> orderings(string& word) {
    string w = word;
    sort(w.begin(), w.end());
    vector<string> out;
    do {
        out.push_back(w);
    } while (next_permutation(w.begin(), w.end()));
    return out;
}
