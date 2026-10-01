#include <queue>
#include <string>
using namespace std;

string firstUnique(string& word) {
    int count[26] = {};
    queue<char> q;
    string out;
    for (char c : word) {
        count[c - 'a']++;
        q.push(c);
        while (!q.empty() && count[q.front() - 'a'] > 1) q.pop();
        out += q.empty() ? '#' : q.front();
    }
    return out;
}
