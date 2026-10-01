#include <algorithm>
#include <string>
#include <vector>
using namespace std;

vector<int> pieceSizes(string& word) {
    int last[26] = {};
    for (int i = 0; i < (int)word.size(); i++) last[word[i] - 'a'] = i;
    vector<int> sizes;
    int start = 0, end = 0;
    for (int i = 0; i < (int)word.size(); i++) {
        end = max(end, last[word[i] - 'a']);
        if (i == end) {
            sizes.push_back(end - start + 1);
            start = i + 1;
        }
    }
    return sizes;
}
